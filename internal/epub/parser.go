package epub

import (
	"archive/zip"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strings"
)

// Book represents a parsed EPUB book
type Book struct {
	Title    string
	Author   string
	Chapters []Chapter
	CSS      []string
	BasePath string
}

// Chapter represents a single chapter/section
type Chapter struct {
	Title   string
	Content string
	Order   int
}

// Container represents the META-INF/container.xml structure
type Container struct {
	XMLName   xml.Name `xml:"container"`
	RootFiles []struct {
		FullPath  string `xml:"full-path,attr"`
		MediaType string `xml:"media-type,attr"`
	} `xml:"rootfiles>rootfile"`
}

// Package represents the OPF package document
type Package struct {
	XMLName  xml.Name `xml:"package"`
	Metadata Metadata `xml:"metadata"`
	Manifest Manifest `xml:"manifest"`
	Spine    Spine    `xml:"spine"`
}

type Metadata struct {
	Title   string `xml:"title"`
	Creator string `xml:"creator"`
}

type Manifest struct {
	Items []ManifestItem `xml:"item"`
}

type ManifestItem struct {
	ID        string `xml:"id,attr"`
	Href      string `xml:"href,attr"`
	MediaType string `xml:"media-type,attr"`
}

type Spine struct {
	ItemRefs []SpineItemRef `xml:"itemref"`
}

type SpineItemRef struct {
	IDRef string `xml:"idref,attr"`
}

// Parse reads and parses an EPUB file
func Parse(epubPath string) (*Book, error) {
	r, err := zip.OpenReader(epubPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open epub: %w", err)
	}
	defer r.Close()

	// Create a map for quick file lookup
	files := make(map[string]*zip.File)
	for _, f := range r.File {
		files[f.Name] = f
	}

	// Parse container.xml to find the OPF file
	containerFile, ok := files["META-INF/container.xml"]
	if !ok {
		return nil, fmt.Errorf("container.xml not found")
	}

	container, err := parseContainer(containerFile)
	if err != nil {
		return nil, err
	}

	if len(container.RootFiles) == 0 {
		return nil, fmt.Errorf("no rootfile found in container.xml")
	}

	opfPath := container.RootFiles[0].FullPath
	opfFile, ok := files[opfPath]
	if !ok {
		return nil, fmt.Errorf("OPF file not found: %s", opfPath)
	}

	pkg, err := parsePackage(opfFile)
	if err != nil {
		return nil, err
	}

	basePath := path.Dir(opfPath)
	if basePath == "." {
		basePath = ""
	}

	book := &Book{
		Title:    pkg.Metadata.Title,
		Author:   pkg.Metadata.Creator,
		BasePath: basePath,
	}

	// Build manifest lookup
	manifestMap := make(map[string]ManifestItem)
	for _, item := range pkg.Manifest.Items {
		manifestMap[item.ID] = item
	}

	// Extract CSS files
	for _, item := range pkg.Manifest.Items {
		if item.MediaType == "text/css" {
			cssPath := resolvePath(basePath, item.Href)
			if cssFile, ok := files[cssPath]; ok {
				content, err := readFileContent(cssFile)
				if err == nil {
					// Embed images referenced in CSS (background-image, etc.)
					cssDir := path.Dir(cssPath)
					content = embedImages(content, cssDir, files)
					book.CSS = append(book.CSS, content)
				}
			}
		}
	}

	// Extract chapters in spine order
	for i, itemRef := range pkg.Spine.ItemRefs {
		item, ok := manifestMap[itemRef.IDRef]
		if !ok {
			continue
		}

		if !strings.Contains(item.MediaType, "html") && !strings.Contains(item.MediaType, "xml") {
			continue
		}

		chapterPath := resolvePath(basePath, item.Href)
		chapterFile, ok := files[chapterPath]
		if !ok {
			continue
		}

		content, err := readFileContent(chapterFile)
		if err != nil {
			continue
		}

		// Process images to embed as base64
		// Use the chapter's directory as the base for resolving relative image paths
		chapterDir := path.Dir(chapterPath)
		content = embedImages(content, chapterDir, files)

		book.Chapters = append(book.Chapters, Chapter{
			Title:   item.ID,
			Content: content,
			Order:   i,
		})
	}

	// Sort chapters by order
	sort.Slice(book.Chapters, func(i, j int) bool {
		return book.Chapters[i].Order < book.Chapters[j].Order
	})

	return book, nil
}

func parseContainer(f *zip.File) (*Container, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var container Container
	if err := xml.NewDecoder(rc).Decode(&container); err != nil {
		return nil, fmt.Errorf("failed to parse container.xml: %w", err)
	}

	return &container, nil
}

func parsePackage(f *zip.File) (*Package, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var pkg Package
	if err := xml.NewDecoder(rc).Decode(&pkg); err != nil {
		return nil, fmt.Errorf("failed to parse OPF: %w", err)
	}

	return &pkg, nil
}

func readFileContent(f *zip.File) (string, error) {
	rc, err := f.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func resolvePath(basePath, href string) string {
	if basePath == "" {
		return href
	}
	return path.Join(basePath, href)
}

func embedImages(content, basePath string, files map[string]*zip.File) string {
	// Regular expressions to find image sources
	// Match src="..." or src='...' in img tags
	imgSrcRegex := regexp.MustCompile(`(?i)(<img[^>]*\ssrc\s*=\s*)(["'])([^"']+)(["'])`)
	
	// Match xlink:href="..." for SVG images
	xlinkRegex := regexp.MustCompile(`(?i)(xlink:href\s*=\s*)(["'])([^"']+)(["'])`)
	
	// Match url(...) in inline styles for background images
	urlRegex := regexp.MustCompile(`(?i)(url\s*\(\s*)(["']?)([^"')]+)(["']?\s*\))`)

	// Process img src attributes
	content = imgSrcRegex.ReplaceAllStringFunc(content, func(match string) string {
		parts := imgSrcRegex.FindStringSubmatch(match)
		if len(parts) < 5 {
			return match
		}
		prefix := parts[1]  // <img ... src=
		quote := parts[2]   // " or '
		src := parts[3]     // the actual path
		endQuote := parts[4]

		dataURI := resolveAndEmbed(src, basePath, files)
		if dataURI != "" {
			return prefix + quote + dataURI + endQuote
		}
		return match
	})

	// Process xlink:href for SVG
	content = xlinkRegex.ReplaceAllStringFunc(content, func(match string) string {
		parts := xlinkRegex.FindStringSubmatch(match)
		if len(parts) < 5 {
			return match
		}
		prefix := parts[1]
		quote := parts[2]
		src := parts[3]
		endQuote := parts[4]

		// Only process image files
		if isImageFile(src) {
			dataURI := resolveAndEmbed(src, basePath, files)
			if dataURI != "" {
				return prefix + quote + dataURI + endQuote
			}
		}
		return match
	})

	// Process url() in styles
	content = urlRegex.ReplaceAllStringFunc(content, func(match string) string {
		parts := urlRegex.FindStringSubmatch(match)
		if len(parts) < 5 {
			return match
		}
		prefix := parts[1]  // url(
		quote := parts[2]   // optional quote
		src := parts[3]     // the path
		suffix := parts[4]  // optional quote + )

		// Only process image files
		if isImageFile(src) {
			dataURI := resolveAndEmbed(src, basePath, files)
			if dataURI != "" {
				return prefix + quote + dataURI + suffix
			}
		}
		return match
	})

	return content
}

func resolveAndEmbed(src, basePath string, files map[string]*zip.File) string {
	// Skip data URIs and external URLs
	if strings.HasPrefix(src, "data:") || strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		return ""
	}

	// Clean up the path (remove fragment identifiers)
	cleanSrc := src
	if idx := strings.Index(cleanSrc, "#"); idx != -1 {
		cleanSrc = cleanSrc[:idx]
	}

	// URL decode the path (handle %20 for spaces, etc.)
	decodedSrc, err := url.QueryUnescape(cleanSrc)
	if err == nil {
		cleanSrc = decodedSrc
	}

	// Resolve relative path
	imagePath := resolveRelativePath(basePath, cleanSrc)

	// Try to find the file
	zipFile, ok := files[imagePath]
	if !ok {
		// Try without basePath
		zipFile, ok = files[cleanSrc]
		if !ok {
			// Try normalizing the path further
			normalizedPath := normalizePath(imagePath)
			zipFile, ok = files[normalizedPath]
			if !ok {
				return ""
			}
		}
	}

	// Read the file content
	data, err := readBinaryContent(zipFile)
	if err != nil {
		return ""
	}

	// Determine MIME type
	mimeType := getMimeType(cleanSrc)
	if mimeType == "" {
		return ""
	}

	// Create data URI
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)
}

func resolveRelativePath(basePath, href string) string {
	if basePath == "" {
		return href
	}

	// Handle ../ in paths
	if strings.HasPrefix(href, "../") {
		// Go up one directory from basePath
		parentDir := path.Dir(basePath)
		return resolveRelativePath(parentDir, strings.TrimPrefix(href, "../"))
	}

	return path.Join(basePath, href)
}

func normalizePath(p string) string {
	// Clean the path and remove leading slashes
	cleaned := path.Clean(p)
	return strings.TrimPrefix(cleaned, "/")
}

func readBinaryContent(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

func isImageFile(src string) bool {
	lower := strings.ToLower(src)
	extensions := []string{".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp", ".bmp", ".ico"}
	for _, ext := range extensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

func getMimeType(filename string) string {
	lower := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(lower, ".jpg"), strings.HasSuffix(lower, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(lower, ".png"):
		return "image/png"
	case strings.HasSuffix(lower, ".gif"):
		return "image/gif"
	case strings.HasSuffix(lower, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(lower, ".webp"):
		return "image/webp"
	case strings.HasSuffix(lower, ".bmp"):
		return "image/bmp"
	case strings.HasSuffix(lower, ".ico"):
		return "image/x-icon"
	default:
		return ""
	}
}

// ToHTML converts the book to a single HTML document
func (b *Book) ToHTML() string {
	var sb strings.Builder

	sb.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	sb.WriteString("<meta charset=\"UTF-8\">\n")
	sb.WriteString(fmt.Sprintf("<title>%s</title>\n", escapeHTML(b.Title)))
	
	// Embed CSS
	sb.WriteString("<style>\n")
	sb.WriteString(`
		body {
			font-family: Georgia, 'Times New Roman', serif;
			line-height: 1.6;
			max-width: 800px;
			margin: 0 auto;
			padding: 40px 20px;
			color: #333;
		}
		h1, h2, h3, h4, h5, h6 {
			margin-top: 1.5em;
			margin-bottom: 0.5em;
		}
		p {
			margin: 0.8em 0;
			text-align: justify;
		}
		img {
			max-width: 100%;
			height: auto;
		}
		.chapter {
			page-break-before: always;
		}
		.chapter:first-child {
			page-break-before: avoid;
		}
		.title-page {
			text-align: center;
			padding: 100px 0;
		}
		.title-page h1 {
			font-size: 2.5em;
			margin-bottom: 0.5em;
		}
		.title-page .author {
			font-size: 1.3em;
			color: #666;
		}
	`)
	for _, css := range b.CSS {
		sb.WriteString(css)
		sb.WriteString("\n")
	}
	sb.WriteString("</style>\n")
	sb.WriteString("</head>\n<body>\n")

	// Title page
	sb.WriteString("<div class=\"title-page\">\n")
	sb.WriteString(fmt.Sprintf("<h1>%s</h1>\n", escapeHTML(b.Title)))
	if b.Author != "" {
		sb.WriteString(fmt.Sprintf("<p class=\"author\">%s</p>\n", escapeHTML(b.Author)))
	}
	sb.WriteString("</div>\n")

	// Chapters
	for _, chapter := range b.Chapters {
		sb.WriteString("<div class=\"chapter\">\n")
		// Extract body content if it's a full HTML document
		content := extractBodyContent(chapter.Content)
		sb.WriteString(content)
		sb.WriteString("\n</div>\n")
	}

	sb.WriteString("</body>\n</html>")

	return sb.String()
}

func extractBodyContent(html string) string {
	// Try to extract just the body content
	bodyStart := strings.Index(strings.ToLower(html), "<body")
	if bodyStart == -1 {
		return html
	}

	// Find the end of the opening body tag
	bodyTagEnd := strings.Index(html[bodyStart:], ">")
	if bodyTagEnd == -1 {
		return html
	}
	bodyStart = bodyStart + bodyTagEnd + 1

	bodyEnd := strings.LastIndex(strings.ToLower(html), "</body>")
	if bodyEnd == -1 {
		bodyEnd = len(html)
	}

	return html[bodyStart:bodyEnd]
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
