
# hugo-docs-ja

[TOC]

## コマンドの実行方法

```ps1
# Go Module の初期化
$ go mod init hugo-docs-i18n

# Cobra CLI アプリケーションのインストール
$ go install github.com/spf13/cobra-cli@latest

# Cobra CLI アプリケーションの起動 (Cobra の初期化)
$ cobra-cli init 
$ go run main.go
$ go run main.go --help

# サブコマンドの追加
$ cobra-cli add version
$ cobra-cli add init
$ cobra-cli add update

# コンパイル
$ go build main.go
$ mv main.exe hugo-docs-i18n.exe
```



-----

## commmands

### 1. import_jekyll.go

Jekyll を Hugo にインポートするコマンド

```go
func convertJekyllPost(path, relPath, targetDir string, draft bool) error {
	jww.TRACE.Println("Converting", path)

	filename := filepath.Base(path)
	postDate, postName, err := parseJekyllFilename(filename)
	if err != nil {
		jww.WARN.Printf("Failed to parse filename '%s': %s. Skipping.", filename, err)
		return nil
	}

	jww.TRACE.Println(filename, postDate, postName)

	targetFile := filepath.Join(targetDir, relPath)
	targetParentDir := filepath.Dir(targetFile)
	os.MkdirAll(targetParentDir, 0777)

	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		jww.ERROR.Println("Read file error:", path)
		return err
	}

	pf, err := pageparser.ParseFrontMatterAndContent(bytes.NewReader(contentBytes))
	if err != nil {
		jww.ERROR.Println("Parse file error:", path)
		return err
	}

	newmetadata, err := convertJekyllMetaData(pf.FrontMatter, postName, postDate, draft)
	if err != nil {
		jww.ERROR.Println("Convert metadata error:", path)
		return err
	}

	content, err := convertJekyllContent(newmetadata, string(pf.Content))
	if err != nil {
		jww.ERROR.Println("Converting Jekyll error:", path)
		return err
	}

	fs := hugofs.Os
	if err := helpers.WriteToDisk(targetFile, strings.NewReader(content), fs); err != nil {
		return fmt.Errorf("failed to save file %q: %s", filename, err)
	}

	return nil
}
```


### 2. convert.go

サイトのコンテンツファイルのフロントマター (JSON/TOML/YAML) を変更するコマンド

```go
func (cc *convertCmd) convertAndSavePage(p page.Page, site *hugolib.Site, targetFormat metadecoders.Format) error {
	// The resources are not in .Site.AllPages.
	for _, r := range p.Resources().ByType("page") {
		if err := cc.convertAndSavePage(r.(page.Page), site, targetFormat); err != nil {
			return err
		}
	}

	if p.File().IsZero() {
		// No content file.
		return nil
	}

	errMsg := fmt.Errorf("Error processing file %q", p.File().Path())
	site.Log.Infoln("Attempting to convert", p.File().Filename())

	f := p.File()
	file, err := f.FileInfo().Meta().Open()
	if err != nil {
		site.Log.Errorln(errMsg)
		file.Close()
		return nil
	}

	pf, err := pageparser.ParseFrontMatterAndContent(file)
	if err != nil {
		site.Log.Errorln(errMsg)
		file.Close()
		return err
	}
	file.Close()

	// better handling of dates in formats that don't have support for them
	if pf.FrontMatterFormat == metadecoders.JSON || pf.FrontMatterFormat == metadecoders.YAML || pf.FrontMatterFormat == metadecoders.TOML {
		for k, v := range pf.FrontMatter {
			switch vv := v.(type) {
			case time.Time:
				pf.FrontMatter[k] = vv.Format(time.RFC3339)
			}
		}
	}

	var newContent bytes.Buffer
	err = parser.InterfaceToFrontMatter(pf.FrontMatter, targetFormat, &newContent)
	if err != nil {
		site.Log.Errorln(errMsg)
		return err
	}
	newContent.Write(pf.Content)
	newFilename := p.File().Filename()

	if cc.outputDir != "" {
		contentDir := strings.TrimSuffix(newFilename, p.File().Path())
		contentDir = filepath.Base(contentDir)
		newFilename = filepath.Join(cc.outputDir, contentDir, p.File().Path())
	}

	fs := hugofs.Os
	if err := helpers.WriteToDisk(newFilename, &newContent, fs); err != nil {
		return fmt.Errorf("Failed to save file %q:: %w", newFilename, err)
	}

	return nil
}
```


### 3. list.go

サイトのコンテンツページのリストを表示するコマンド

```go
		&cobra.Command{
			Use:   "all",
			Short: "List all posts",
			Long:  `List all of the posts in your content directory, include drafts, future and expired pages.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				sites, err := cc.buildSites(map[string]any{
					"buildExpired": true,
					"buildDrafts":  true,
					"buildFuture":  true,
				})
				if err != nil {
					return newSystemError("Error building sites", err)
				}

				writer := csv.NewWriter(os.Stdout)
				defer writer.Flush()

				writer.Write([]string{
					"path",
					"slug",
					"title",
					"date",
					"expiryDate",
					"publishDate",
					"draft",
					"permalink",
				})
				for _, p := range sites.Pages() {
					if !p.IsPage() {
						continue
					}
					err := writer.Write([]string{
						strings.TrimPrefix(p.File().Filename(), sites.WorkingDir+string(os.PathSeparator)),
                                                 // PageWithoutContent.FileProvider
						p.Slug(),                            // PageMetaProvider
						p.Title(),                           // PageWithoutContent.Resource.ResourceMetaProvider
						p.Date().Format(time.RFC3339),       // PageMetaProvider.Dated
						p.ExpiryDate().Format(time.RFC3339),  //
						p.PublishDate().Format(time.RFC3339), //
						strconv.FormatBool(p.Draft()),        // PageMetaProvider
						p.Permalink(),                        // PageWithoutContent.Resource.ResourceLinksProvider
					})
					if err != nil {
						return newSystemError("Error writing posts to stdout", err)
					}
				}

				return nil
			},
		},
```




-----

## hugolib

### hugolib/page_meta.go

```go
type pageMeta struct {
}

func (p *pageMeta) Lang() string {
	return p.s.Lang()
}

func (p *pageMeta) Draft() bool {
	return p.draft
}

func (p *pageMeta) File() source.File {
	return p.f
}

func (pm *pageMeta) setMetadata(parentBucket *pagesMapBucket, p *pageState, frontmatter map[string]any) error {
  // マップ frontmatter を struct pegeMeta の各フィールドに設定する
}
```


### hugolib/page.go

```go
type pageState struct {
	// This slice will be of same length as the number of global slice of output
	// formats (for all sites).
	pageOutputs []*pageOutput

	// Used to determine if we can reuse content across output formats.
	pageOutputTemplateVariationsState *atomic.Uint32

	// This will be shifted out when we start to render a new output format.
	*pageOutput

	// Common for all output formats.
	*pageCommon
}

func (p *pageState) MarshalJSON() ([]byte, error) {
	return page.MarshalPageToJSON(p)
}

func (p *pageState) getPages() page.Pages {
	b := p.bucket
	if b == nil {
		return nil
	}
	return b.getPages()
}

func (p *pageState) Pages() page.Pages {
	p.pagesInit.Do(func() {
		var pages page.Pages

		switch p.Kind() {
		case page.KindPage:
		case page.KindSection, page.KindHome:
			pages = p.getPagesAndSections()
		case page.KindTerm:
			pages = p.bucket.getTaxonomyEntries()
		case page.KindTaxonomy:
			pages = p.bucket.getTaxonomies()
		default:
			pages = p.s.Pages()
		}
		p.pages = pages
	})
	return p.pages
}

func (p *pageState) mapContentForResult(
	result pageparser.Result,
	s *shortcodeHandler,
	rn *pageContentMap,
	markup string,
	withFrontMatter func(map[string]any) error,
) error {
  // ...
Loop:
	for {
		it := iter.Next()

		switch {
		case it.Type == pageparser.TypeIgnore:
		case it.IsFrontMatter():
			f := pageparser.FormatFromFrontMatterType(it.Type)
			m, err := metadecoders.Default.UnmarshalToMap(it.Val, f)
    }
  }
}
```


### hugolib/hugo_sites_build.go

```go
// Build builds all sites. If filesystem events are provided,
// this is considered to be a potential partial rebuild.
func (h *HugoSites) Build(config BuildCfg, events ...fsnotify.Event) error {

}
```


### hugolib/hugo_sites.go

```go
// HugoSites represents the sites to build. Each site represents a language.
type HugoSites struct {
	Sites []*Site
	multilingual *Multilingual
	// Multihost is set if multilingual and baseURL set on the language level.
	multihost bool
	// If this is running in the dev server.
	running bool
	// Render output formats for all sites.
	renderFormats output.Formats
	// The currently rendered Site.
	currentSite *Site
	*deps.Deps
	gitInfo       *gitInfo
	codeownerInfo *codeownerInfo
	// As loaded from the /data dirs
	data map[string]any
	contentInit sync.Once
	content     *pageMaps
	// Keeps track of bundle directories and symlinks to enable partial rebuilding.
	ContentChanges *contentChangeMap
	// File change events with filename stored in this map will be skipped.
	skipRebuildForFilenamesMu sync.Mutex
	skipRebuildForFilenames   map[string]bool
	init *hugoSitesInit
	workers    *para.Workers
	numWorkers int
	*fatalErrorHandler
	*testCounters
}

// NewHugoSites creates HugoSites from the given config.
func NewHugoSites(cfg deps.DepsCfg) (*HugoSites, error) {
	if cfg.Logger == nil {
		cfg.Logger = loggers.NewErrorLogger()
	}
	sites, err := createSitesFromConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("from config: %w", err)
	}
	return newHugoSites(cfg, sites...)
}

// NewHugoSites creates a new collection of sites given the input sites, building
// a language configuration based on those.
func newHugoSites(cfg deps.DepsCfg, sites ...*Site) (*HugoSites, error) {
  // ...
  h := &HugoSites{
		running:                 cfg.Running,
		multilingual:            langConfig,
		multihost:               cfg.Cfg.GetBool("multihost"),
		Sites:                   sites,
		workers:                 workers,
		numWorkers:              numWorkers,
		skipRebuildForFilenames: make(map[string]bool),
		init: &hugoSitesInit{
			data:         lazy.New(),
			layouts:      lazy.New(),
			gitInfo:      lazy.New(),
			translations: lazy.New(),
		},
	}
  // ...
	for _, s := range sites {
		s.h = h
	}
  // ...
 	return h, initErr
}

// Config に書かれた複数の言語ごとにサイトを作成し、サイトの配列を返す
func createSitesFromConfig(cfg deps.DepsCfg) ([]*Site, error) {
	var sites []*Site
	languages := getLanguages(cfg.Cfg)
	for _, lang := range languages {
		if lang.Disabled {
			continue
		}
		var s *Site
		var err error
		cfg.Language = lang
		s, err = newSite(cfg)
		if err != nil {
			return nil, err
		}
		sites = append(sites, s)
	}
	return sites, nil
}
```


### hugolib/site.go

```go
// サイトには、静的サイトを構築するために必要な情報がすべて含まれています。
// 基本的な情報の流れは以下の通りです。
//
// 1. Files のリストを解析し、Pages に変換します。
// 2. Pages にはセクション (生成元のファイルに基づく)、エイリアス、スラッグ (ページのフロントマターに含まれる) が含まれ、
//    これらは生成されるさまざまなターゲットとなります。 
//    それらは標準的なリストとなります。 標準的なパスは、パターンに基づいて上書きすることができます。
// 3. タクソノミーは設定によって作成され、最終ページの何らかの側面と、通常はパーマ URL を提示することになります。
// 4. すべての Pages は、多数の異なる要素に基づいた目的のレイアウトに基づき、テンプレートを通過します。
// 5. Files のコレクション全体がディスクに書き込まれます。

type Site struct {
	// 所有するコンテナ。 多言語の場合、複数のサイトが存在することになる。
	h *HugoSites
	*PageCollections
  taxonomies TaxonomyList
	Sections Taxonomy
	Info     *SiteInfo
	language   *langs.Language
	siteBucket *pagesMapBucket
  // ...
 	// このサイトの最終更新日
	lastmod time.Time
  // ...
}

func (s *Site) Language() *langs.Language {
	return s.language
}

// newSite creates a new site with the given configuration.
func newSite(cfg deps.DepsCfg) (*Site, error) {
  // ...
	titleFunc := helpers.GetTitleFunc(cfg.Language.GetString("titleCaseStyle"))
	frontMatterHandler, err := pagemeta.NewFrontmatterHandler(cfg.Logger, cfg.Cfg)
  // ...
	s := &Site{
		language:      cfg.Language,
		siteBucket:    siteBucket,
		disabledKinds: disabledKinds,
		outputFormats:       outputFormats,
		outputFormatsConfig: siteOutputFormatsConfig,
		mediaTypesConfig:    siteMediaTypesConfig,
		siteCfg: siteConfig,
		titleFunc: titleFunc,
		rc: &siteRenderingContext{output.HTMLFormat},
		frontmatterHandler: frontMatterHandler,
		relatedDocsHandler: page.NewRelatedDocsHandler(relatedContentConfig),
	}
	s.prepareInits()
	return s, nil
}

// NewSite creates a new site with the given dependency configuration.
// The site will have a template system loaded and ready to use.
// Note: This is mainly used in single site tests.
func NewSite(cfg deps.DepsCfg) (*Site, error) {
	s, err := newSite(cfg)
  // ...
}
```


-----

## parser

### parser/pageparser/pageparser.go

```go
type ContentFrontMatter struct {
	Content           []byte
	FrontMatter       map[string]any
	FrontMatterFormat metadecoders.Format
}

// ParseFrontMatterAndContent is a convenience method to extract front matter
// and content from a content page.
func ParseFrontMatterAndContent(r io.Reader) (ContentFrontMatter, error) {
	var cf ContentFrontMatter

}

func FormatFromFrontMatterType(typ ItemType) metadecoders.Format {
	switch typ {
	case TypeFrontMatterJSON:
		return metadecoders.JSON
	case TypeFrontMatterORG:
		return metadecoders.ORG
	case TypeFrontMatterTOML:
		return metadecoders.TOML
	case TypeFrontMatterYAML:
		return metadecoders.YAML
	default:
		return ""
	}
}
```


### parser/frontmatter.go

```go
func InterfaceToConfig(in any, format metadecoders.Format, w io.Writer) error {
}

func InterfaceToFrontMatter(in any, format metadecoders.Format, w io.Writer) error {
}
```


### parser/metadecoder/decoder.go

```go
type Decoder struct {
	// Delimiter is the field delimiter used in the CSV decoder. It defaults to ','.
	Delimiter rune

	// Comment, if not 0, is the comment character ued in the CSV decoder. Lines beginning with the
	// Comment character without preceding whitespace are ignored.
	Comment rune
}

// Default is a Decoder in its default configuration.
var Default = Decoder{
	Delimiter: ',',
}

// UnmarshalToMap will unmarshall data in format f into a new map. This is
// what's needed for Hugo's front matter decoding.
func (d Decoder) UnmarshalToMap(data []byte, f Format) (map[string]any, error) {
	m := make(map[string]any)
  // ...
  err := d.UnmarshalTo(data, f, &m)
	return m, err
}

// UnmarshalFileToMap is the same as UnmarshalToMap, but reads the data from
// the given filename.
func (d Decoder) UnmarshalFileToMap(fs afero.Fs, filename string) (map[string]any, error) {
	format := FormatFromString(filename)
  // ...
  return d.UnmarshalToMap(data, format)
}

// UnmarshalStringTo tries to unmarshal data to a new instance of type typ.
func (d Decoder) UnmarshalStringTo(data string, typ any) (any, error) {
	data = strings.TrimSpace(data)

}

// Unmarshal will unmarshall data in format f into an interface{}.
// This is what's needed for Hugo's /data handling.
func (d Decoder) Unmarshal(data []byte, f Format) (any, error) {

}

// UnmarshalTo unmarshals data in format f into v.
func (d Decoder) UnmarshalTo(data []byte, f Format, v any) error {
	var err error

	switch f {
	case ORG:
		err = d.unmarshalORG(data, v)
	case JSON:
		err = json.Unmarshal(data, v)
	case XML:
		var xmlRoot xml.Map
		xmlRoot, err = xml.NewMapXml(data)

		var xmlValue map[string]any
		if err == nil {
			xmlRootName, err := xmlRoot.Root()
			if err != nil {
				return toFileError(f, data, fmt.Errorf("failed to unmarshal XML: %w", err))
			}
			xmlValue = xmlRoot[xmlRootName].(map[string]any)
		}

		switch v := v.(type) {
		case *map[string]any:
			*v = xmlValue
		case *any:
			*v = xmlValue
		}
	case TOML:
		err = toml.Unmarshal(data, v)
	case YAML:
		err = yaml.Unmarshal(data, v)
		if err != nil {
			return toFileError(f, data, fmt.Errorf("failed to unmarshal YAML: %w", err))
		}

		// To support boolean keys, the YAML package unmarshals maps to
		// map[interface{}]interface{}. Here we recurse through the result
		// and change all maps to map[string]interface{} like we would've
		// gotten from `json`.
		var ptr any
		switch v.(type) {
		case *map[string]any:
			ptr = *v.(*map[string]any)
		case *any:
			ptr = *v.(*any)
		default:
			// Not a map.
		}

		if ptr != nil {
			if mm, changed := stringifyMapKeys(ptr); changed {
				switch v.(type) {
				case *map[string]any:
					*v.(*map[string]any) = mm.(map[string]any)
				case *any:
					*v.(*any) = mm
				}
			}
		}
	case CSV:
		return d.unmarshalCSV(data, v)

	default:
		return fmt.Errorf("unmarshal of format %q is not supported", f)
	}

	if err == nil {
		return nil
	}

	return toFileError(f, data, fmt.Errorf("unmarshal failed: %w", err))
}
```


### parser/metadecoder/format.go

```go
const (
	// These are the supported metdata  formats in Hugo. Most of these are also
	// supported as /data formats.
	ORG  Format = "org"
	JSON Format = "json"
	TOML Format = "toml"
	YAML Format = "yaml"
	CSV  Format = "csv"
	XML  Format = "xml"
)

// FormatFromString turns formatStr, typically a file extension without any ".",
// into a Format. It returns an empty string for unknown formats.
func FormatFromString(formatStr string) Format {

}
```



-----

## resources

### resources/pager/page.go

```go
// Page is the core interface in Hugo.
type Page interface {
	ContentProvider           // 1.
	TableOfContentsProvider   // 2.
	PageWithoutContent        // 3.
}

// 1. ContentProvider provides the content related values for a Page.
type ContentProvider interface {
	Content() (any, error)
	// Plain returns the Page Content stripped of HTML markup.
	Plain() string
	// PlainWords returns a string slice from splitting Plain using https://pkg.go.dev/strings#Fields.
	PlainWords() []string
	// Summary returns a generated summary of the content.
	// The breakpoint can be set manually by inserting a summary separator in the source file.
	Summary() template.HTML
	// Truncated returns whether the Summary  is truncated or not.
	Truncated() bool
	// FuzzyWordCount returns the approximate number of words in the content.
	FuzzyWordCount() int
	// WordCount returns the number of words in the content.
	WordCount() int
	// ReadingTime returns the reading time based on the length of plain text.
	ReadingTime() int
	// Len returns the length of the content.
	Len() int
}

// 2. TableOfContentsProvider provides the table of contents for a Page.
type TableOfContentsProvider interface {
	TableOfContents() template.HTML
}

// 3. PageWithoutContent is the Page without any of the content methods.
type PageWithoutContent interface {
	RawContentProvider        // 4. 
	resource.Resource         // 13.
	PageMetaProvider          // 5.
	resource.LanguageProvider // 15. 
	// For pages backed by a file.
	FileProvider              // 6.
	GitInfoProvider
	// Output formats
	OutputFormatsProvider     // 7.
	AlternativeOutputFormatsProvider
	// Tree navigation
	ChildCareProvider
	TreeProvider
	// Horizontal navigation
	InSectionPositioner
	PageRenderProvider
	PaginatorProvider
	Positioner
	navigation.PageMenusProvider
	// TODO(bep)
	AuthorProvider
	// Page lookups/refs
	GetPageProvider
	RefProvider
	resource.TranslationKeyProvider
	TranslationsProvider
	SitesProvider
	// Helper methods
	ShortcodeInfoProvider
	compare.Eqer
	// Scratch returns a Scratch that can be used to store temporary state.
	// Note that this Scratch gets reset on server rebuilds. See Store() for a variant that survives.
	maps.Scratcher
	// Store returns a Scratch that can be used to store temporary state.
	// In contrast to Scratch(), this Scratch is not reset on server rebuilds.
	Store() *maps.Scratch
	RelatedKeywordsProvider
	// GetTerms gets the terms of a given taxonomy,
	// e.g. GetTerms("categories")
	GetTerms(taxonomy string) Pages
	// Used in change/dependency tracking.
	identity.Provider
	DeprecatedWarningPageMethods
}

// 4. RawContentProvider provides the raw, unprocessed content of the page.
type RawContentProvider interface {
	RawContent() string
}

// 5. PageMetaProvider provides page metadata, typically provided via front matter.
type PageMetaProvider interface {
	// The 4 page dates
	resource.Dated                          // 11. hugo/resources/resource/dates.go
	// Aliases forms the base for redirects generation.
	Aliases() []string
	// BundleType returns the bundle type: `leaf`, `branch` or an empty string.
	BundleType() files.ContentClass
	// A configured description.
	Description() string
	// Whether this is a draft. Will only be true if run with the --buildDrafts (-D) flag.
	Draft() bool
	// IsHome returns whether this is the home page.
	IsHome() bool
	// Configured keywords.
	Keywords() []string
	// The Page Kind. One of page, home, section, taxonomy, term.
	Kind() string
	// The configured layout to use to render this page. Typically set in front matter.
	Layout() string
	// The title used for links.
	LinkTitle() string
	// IsNode returns whether this is an item of one of the list types in Hugo,
	// i.e. not a regular content
	IsNode() bool
	// IsPage returns whether this is a regular content
	IsPage() bool
	// Param looks for a param in Page and then in Site config.
	Param(key any) (any, error)
	// Path gets the relative path, including file name and extension if relevant,
	// to the source of this Page. It will be relative to any content root.
	Path() string
	// This is just a temporary bridge method. Use Path in templates.
	// Pathc is for internal usage only.
	Pathc() string
	// The slug, typically defined in front matter.
	Slug() string
	// This page's language code. Will be the same as the site's.
	Lang() string
	// IsSection returns whether this is a section
	IsSection() bool
	// Section returns the first path element below the content root.
	Section() string
	// Returns a slice of sections (directories if it's a file) to this
	// Page.
	SectionsEntries() []string
	// SectionsPath is SectionsEntries joined with a /.
	SectionsPath() string
	// Sitemap returns the sitemap configuration for this page.
	Sitemap() config.Sitemap
	// Type is a discriminator used to select layouts etc. It is typically set
	// in front matter, but will fall back to the root section.
	Type() string
	// The configured weight, used as the first sort value in the default
	// page sort if non-zero.
	Weight() int
}

// 6. FileProvider provides the source file.
type FileProvider interface {
	File() source.File
}

// 7. OutputFormatsProvider provides the OutputFormats of a Page.
type OutputFormatsProvider interface {
	OutputFormats() OutputFormats             // 12. 
}
```


### resources/page/pages.go

```go
// Pages is a slice of Page objects. This is the most common list type in Hugo.
type Pages []Page
```


### resources/resource/dates.go

```go
// 11. Dated wraps a "dated resource". These are the 4 dates that makes
// the date logic in Hugo.
type Dated interface {
	// Date returns the date of the resource.
	Date() time.Time
	// Lastmod returns the last modification date of the resource.
	Lastmod() time.Time
	// PublishDate returns the publish date of the resource.
	PublishDate() time.Time
	// ExpiryDate returns the expiration date of the resource.
	ExpiryDate() time.Time
}

// Dates holds the 4 Hugo dates.
type Dates struct {
	FDate        time.Time
	FLastmod     time.Time
	FPublishDate time.Time
	FExpiryDate  time.Time
}

func (p Dates) Date() time.Time {
	return p.FDate
}

func (p Dates) Lastmod() time.Time {
	return p.FLastmod
}

func (p Dates) PublishDate() time.Time {
	return p.FPublishDate
}

func (p Dates) ExpiryDate() time.Time {
	return p.FExpiryDate
}
```


### resources/page/page_outputformat.go

```go
// OutputFormats holds a list of the relevant output formats for a given page.
type OutputFormats []OutputFormat

// 12. OutputFormat links to a representation of a resource.
type OutputFormat struct {
	// Rel contains a value that can be used to construct a rel link.
	// This is value is fetched from the output format definition.
	// Note that for pages with only one output format,
	// this method will always return "canonical".
	// As an example, the AMP output format will, by default, return "amphtml".
	//
	// See:
	// https://www.ampproject.org/docs/guides/deploy/discovery
	//
	// Most other output formats will have "alternate" as value for this.
	Rel string
	Format output.Format
	relPermalink string
	permalink    string
}

// Permalink returns the absolute permalink to this output format.
func (o OutputFormat) Permalink() string {
	return o.permalink
}
```


### resources/resource/resourcetypes.go

```go
// 13. Resource represents a linkable resource, i.e. a content page, image etc.
type Resource interface {
	ResourceTypeProvider
	MediaTypeProvider
	ResourceLinksProvider     // 14. 
	ResourceMetaProvider      // 15. 
	ResourceParamsProvider
	ResourceDataProvider
	ErrProvider
}

// 14. 
type ResourceLinksProvider interface {
	// Permalink represents the absolute link to this resource.
	Permalink() string
	// RelPermalink represents the host relative link to this resource.
	RelPermalink() string
}

// 15. 
type ResourceMetaProvider interface {
	// Name is the logical name of this resource. This can be set in the front matter
	// metadata for this resource. If not set, Hugo will assign a value.
	// This will in most cases be the base filename.
	// So, for the image "/some/path/sunset.jpg" this will be "sunset.jpg".
	// The value returned by this method will be used in the GetByPrefix and ByPrefix methods
	// on Resources.
	Name() string
	// Title returns the title if set in front matter. For content pages, this will be the expected value.
	Title() string
}

// 15. LanguageProvider is a Resource in a language.
type LanguageProvider interface {
	Language() *langs.Language
}
```



-----

## langs

### langs/language.go

```go

// Language manages specific-language configuration.
type Language struct {
	Lang              string
	LanguageName      string
	LanguageDirection string
	Title             string
	Weight            int
	// For internal use.
	Disabled bool
	// If set per language, this tells Hugo that all content files without any
	// language indicator (e.g. my-page.en.md) is in this language.
	// This is usually a path relative to the working dir, but it can be an
	// absolute directory reference. It is what we get.
	// For internal use.
	ContentDir string
	// Global config.
	// For internal use.
	Cfg config.Provider
	// Language specific config.
	// For internal use.
	LocalCfg config.Provider
	// Composite config.
	// For internal use.
	config.Provider
	// These are params declared in the [params] section of the language merged with the
	// site's params, the most specific (language) wins on duplicate keys.
	params    map[string]any
	paramsMu  sync.Mutex
	paramsSet bool
	// Used for date formatting etc. We don't want these exported to the
	// templates.
	// TODO(bep) do the same for some of the others.
	translator    locales.Translator
	timeFormatter htime.TimeFormatter
	tag           language.Tag
	collator      *Collator
	location      *time.Location
	// Error during initialization. Will fail the buld.
	initErr error
}

// For internal use.
func (l *Language) String() string {
	return l.Lang
}

// NewLanguage creates a new language.
func NewLanguage(lang string, cfg config.Provider) *Language {
  // ...
	l := &Language{
		Lang:       lang,
		ContentDir: cfg.GetString("contentDir"),
		Cfg:        cfg, 
    LocalCfg:   localCfg,
		Provider:      compositeConfig,
		params:        params,
		translator:    translator,
		timeFormatter: htime.NewTimeFormatter(translator),
		tag:           tag,
		collator:      coll,
	}
	if err := l.loadLocation(cfg.GetString("timeZone")); err != nil {
		l.initErr = err
	}
	return l
}

// IsMultihost returns whether there are more than one language and at least one of
// the languages has baseURL specificed on the language level.
func (l Languages) IsMultihost() bool {
}

// GetLocal gets a configuration value set on language level. It will
// not fall back to any global value.
// It will return nil if a value with the given key cannot be found.
// For internal use.
func (l *Language) GetLocal(key string) any {
}
```
