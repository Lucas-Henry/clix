package icons

import "strings"

const (
	IconDir        = "\uf115" // nf-fa-folder_open
	IconDirEmpty   = "\uf114" // nf-fa-folder
	IconSymlink    = "\uf481" // nf-oct-file_symlink_file
	IconFile       = "\uf15b" // nf-fa-file
	IconGo         = "\ue627" // nf-dev-go
	IconRust       = "\ue7a8" // nf-dev-rust
	IconPython     = "\ue606" // nf-dev-python
	IconJS         = "\ue74e" // nf-dev-javascript
	IconTS         = "\ue628" // nf-dev-typescript
	IconHTML       = "\ue736" // nf-dev-html5
	IconCSS        = "\ue749" // nf-dev-css3
	IconJSON       = "\ue60b" // nf-seti-json
	IconYAML       = "\ue6a8" // nf-seti-yml
	IconTOML       = "\ue6b2" // nf-seti-config
	IconMarkdown   = "\uf48a" // nf-oct-markdown
	IconShell      = "\uf489" // nf-dev-terminal
	IconGit        = "\ue702" // nf-dev-git
	IconDockerfile = "\uf308" // nf-linux-docker
	IconImage      = "\uf1c5" // nf-fa-file_image_o
	IconVideo      = "\uf03d" // nf-fa-video_camera
	IconAudio      = "\uf001" // nf-fa-music
	IconArchive    = "\uf410" // nf-oct-file_zip
	IconPDF        = "\uf1c1" // nf-fa-file_pdf_o
	IconLock       = "\uf023" // nf-fa-lock
	IconConfig     = "\ue615" // nf-seti-config
	IconDatabase   = "\uf1c0" // nf-fa-database
	IconC          = "\ue61e" // nf-custom-c
	IconCpp        = "\ue61d" // nf-custom-cpp
	IconJava       = "\ue738" // nf-dev-java
	IconRuby       = "\ue739" // nf-dev-ruby
	IconPHP        = "\ue73d" // nf-dev-php
	IconSwift      = "\ue755" // nf-dev-swift
	IconKotlin     = "\ue634" // nf-dev-kotlin
	IconLua        = "\ue620" // nf-seti-lua
	IconVue        = "\uf0844" // nf-md-vuejs
	IconReact      = "\ue7ba" // nf-dev-react
	IconSvelte     = "\ue697" // nf-seti-svelte
	IconNix        = "\uf313" // nf-linux-nixos
	IconMakefile   = "\ue673" // nf-seti-makefile
	IconEnv        = "\uf462" // nf-seti-dotenv (hidden/env)
	IconBinary     = "\uf471" // nf-oct-file_binary
	IconXML        = "\uf05c0" // nf-md-xml
	IconCSV        = "\uf1c3" // nf-fa-file_excel_o
	IconSQL        = "\uf1c0" // nf-fa-database
	IconProto      = "\uf7ef" // nf-mdi-protocol
)

var specialNames = map[string]string{
	"dockerfile":      IconDockerfile,
	"docker-compose":  IconDockerfile,
	"makefile":        IconMakefile,
	"gnumakefile":     IconMakefile,
	".gitignore":      IconGit,
	".gitmodules":     IconGit,
	".gitattributes":  IconGit,
	".gitconfig":      IconGit,
	".env":            IconEnv,
	".env.local":      IconEnv,
	".env.example":    IconEnv,
	".npmrc":          IconConfig,
	".nvmrc":          IconConfig,
	".editorconfig":   IconConfig,
	"package.json":    IconJS,
	"package-lock.json": IconJS,
	"yarn.lock":       IconJS,
	"cargo.toml":      IconRust,
	"cargo.lock":      IconRust,
	"go.mod":          IconGo,
	"go.sum":          IconGo,
	"procfile":        IconConfig,
	"readme.md":       IconMarkdown,
	"license":         IconLock,
	"licence":         IconLock,
}

var extIcons = map[string]string{
	"go":      IconGo,
	"rs":      IconRust,
	"py":      IconPython,
	"js":      IconJS,
	"mjs":     IconJS,
	"cjs":     IconJS,
	"jsx":     IconReact,
	"ts":      IconTS,
	"tsx":     IconReact,
	"vue":     IconVue,
	"svelte":  IconSvelte,
	"html":    IconHTML,
	"htm":     IconHTML,
	"css":     IconCSS,
	"scss":    IconCSS,
	"sass":    IconCSS,
	"less":    IconCSS,
	"json":    IconJSON,
	"jsonc":   IconJSON,
	"yaml":    IconYAML,
	"yml":     IconYAML,
	"toml":    IconTOML,
	"ini":     IconConfig,
	"conf":    IconConfig,
	"config":  IconConfig,
	"md":      IconMarkdown,
	"mdx":     IconMarkdown,
	"rst":     IconMarkdown,
	"sh":      IconShell,
	"bash":    IconShell,
	"zsh":     IconShell,
	"fish":    IconShell,
	"nu":      IconShell,
	"ps1":     IconShell,
	"c":       IconC,
	"h":       IconC,
	"cpp":     IconCpp,
	"cc":      IconCpp,
	"cxx":     IconCpp,
	"hpp":     IconCpp,
	"java":    IconJava,
	"kt":      IconKotlin,
	"kts":     IconKotlin,
	"rb":      IconRuby,
	"erb":     IconRuby,
	"php":     IconPHP,
	"swift":   IconSwift,
	"lua":     IconLua,
	"xml":     IconXML,
	"svg":     IconImage,
	"png":     IconImage,
	"jpg":     IconImage,
	"jpeg":    IconImage,
	"gif":     IconImage,
	"webp":    IconImage,
	"ico":     IconImage,
	"bmp":     IconImage,
	"mp4":     IconVideo,
	"mkv":     IconVideo,
	"avi":     IconVideo,
	"mov":     IconVideo,
	"webm":    IconVideo,
	"mp3":     IconAudio,
	"flac":    IconAudio,
	"wav":     IconAudio,
	"ogg":     IconAudio,
	"m4a":     IconAudio,
	"zip":     IconArchive,
	"tar":     IconArchive,
	"gz":      IconArchive,
	"bz2":     IconArchive,
	"xz":      IconArchive,
	"zst":     IconArchive,
	"7z":      IconArchive,
	"rar":     IconArchive,
	"pdf":     IconPDF,
	"sql":     IconSQL,
	"db":      IconDatabase,
	"sqlite":  IconDatabase,
	"sqlite3": IconDatabase,
	"csv":     IconCSV,
	"tsv":     IconCSV,
	"proto":   IconProto,
	"nix":     IconNix,
	"lock":    IconLock,
	"env":     IconEnv,
	"bin":     IconBinary,
	"exe":     IconBinary,
	"so":      IconBinary,
	"dylib":   IconBinary,
	"out":     IconBinary,
}

func ForEntry(name string, isDir bool, isEmpty bool) string {
	if isDir {
		if isEmpty {
			return IconDirEmpty
		}
		return IconDir
	}

	lower := strings.ToLower(name)

	if icon, ok := specialNames[lower]; ok {
		return icon
	}

	if idx := strings.LastIndex(lower, "."); idx >= 0 {
		ext := lower[idx+1:]
		if icon, ok := extIcons[ext]; ok {
			return icon
		}
	}

	return IconFile
}
