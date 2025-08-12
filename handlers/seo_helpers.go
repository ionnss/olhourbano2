package handlers

// Breadcrumb represents a single breadcrumb item
type Breadcrumb struct {
	Title    string
	URL      string
	IsActive bool
}

// SEOData contains SEO-related data for pages
type SEOData struct {
	Title       string
	Description string
	Keywords    string
	Canonical   string
	Breadcrumbs []Breadcrumb
	PageSlug    string
}

// GenerateBreadcrumbs creates breadcrumb data for different page types
func GenerateBreadcrumbs(pageType string, additionalData map[string]string) []Breadcrumb {
	baseBreadcrumbs := []Breadcrumb{
		{Title: "Início", URL: "/", IsActive: false},
	}

	switch pageType {
	case "feed":
		return append(baseBreadcrumbs, Breadcrumb{
			Title:    "Denúncias Recentes",
			URL:      "/feed",
			IsActive: true,
		})
	case "map":
		return append(baseBreadcrumbs, Breadcrumb{
			Title:    "Mapa de Denúncias",
			URL:      "/map",
			IsActive: true,
		})
	case "report":
		return append(baseBreadcrumbs, Breadcrumb{
			Title:    "Fazer Denúncia",
			URL:      "/report",
			IsActive: true,
		})
	case "report_detail":
		reportID := additionalData["reportID"]
		return append(baseBreadcrumbs,
			Breadcrumb{Title: "Denúncias Recentes", URL: "/feed", IsActive: false},
			Breadcrumb{Title: "Denúncia #" + reportID, URL: "/report/" + reportID, IsActive: true},
		)
	case "footer_page":
		pageTitle := additionalData["pageTitle"]
		pageSlug := additionalData["pageSlug"]
		return append(baseBreadcrumbs, Breadcrumb{
			Title:    pageTitle,
			URL:      "/" + pageSlug,
			IsActive: true,
		})
	default:
		return baseBreadcrumbs
	}
}

// GenerateSEOData creates comprehensive SEO data for pages
func GenerateSEOData(pageType string, additionalData map[string]string) SEOData {
	baseURL := "https://olhourbano.com.br"

	switch pageType {
	case "index":
		return SEOData{
			Title:       "Olho Urbano - Plataforma de Denúncias Urbanas | Reporte Problemas na Cidade",
			Description: "Reporte problemas urbanos e acompanhe denúncias em tempo real. Plataforma cidadã para melhorar sua cidade. Denuncie buracos, iluminação, segurança e mais.",
			Keywords:    "denúncia pública, infraestrutura urbana, problemas na cidade, buracos na rua, iluminação pública, segurança pública, cidadania ativa, corrupção, meio ambiente, obras, fiscalização de comércio, denúncia online, problemas urbanos",
			Canonical:   baseURL + "/",
			Breadcrumbs: GenerateBreadcrumbs("index", nil),
		}
	case "feed":
		return SEOData{
			Title:       "Denúncias Recentes - Olho Urbano | Acompanhe Problemas Urbanos",
			Description: "Acompanhe denúncias recentes de problemas urbanos. Visualize status, fotos e comentários de denúncias cidadãs. Transparência em tempo real.",
			Keywords:    "denúncias recentes, problemas urbanos, acompanhar denúncias, status denúncias, transparência pública, denúncias cidadãs, problemas na cidade, buracos, iluminação, segurança",
			Canonical:   baseURL + "/feed",
			Breadcrumbs: GenerateBreadcrumbs("feed", nil),
		}
	case "map":
		return SEOData{
			Title:       "Mapa de Denúncias - Olho Urbano | Visualize Problemas na Cidade",
			Description: "Visualize denúncias urbanas em um mapa interativo. Encontre problemas na sua região, acompanhe status e localização das denúncias cidadãs.",
			Keywords:    "mapa de denúncias, problemas urbanos mapa, denúncias por localização, mapa interativo, problemas na cidade, localização denúncias, mapa urbano",
			Canonical:   baseURL + "/map",
			Breadcrumbs: GenerateBreadcrumbs("map", nil),
		}
	case "report":
		return SEOData{
			Title:       "Fazer Denúncia - Olho Urbano | Reporte Problemas Urbanos",
			Description: "Faça sua denúncia de forma simples e rápida. Reporte problemas de infraestrutura, segurança, meio ambiente e mais. Plataforma cidadã para melhorar sua cidade.",
			Keywords:    "fazer denúncia, denúncia online, reportar problema, infraestrutura urbana, buracos na rua, iluminação pública, segurança pública, meio ambiente, obras, fiscalização, denúncia cidadã",
			Canonical:   baseURL + "/report",
			Breadcrumbs: GenerateBreadcrumbs("report", nil),
		}
	case "report_detail":
		reportID := additionalData["reportID"]
		pageTitle := additionalData["pageTitle"]
		return SEOData{
			Title:       pageTitle + " - Olho Urbano",
			Description: "Visualize detalhes da denúncia #" + reportID + " no Olho Urbano. Acompanhe status, fotos, comentários e atualizações desta denúncia urbana.",
			Keywords:    "denúncia #" + reportID + ", detalhes denúncia, status denúncia, problemas urbanos, denúncia cidadã, acompanhar denúncia, fotos denúncia, comentários denúncia",
			Canonical:   baseURL + "/report/" + reportID,
			Breadcrumbs: GenerateBreadcrumbs("report_detail", additionalData),
		}
	case "footer_page":
		pageTitle := additionalData["pageTitle"]
		pageSubtitle := additionalData["pageSubtitle"]
		pageSlug := additionalData["pageSlug"]
		return SEOData{
			Title:       pageTitle + " - Olho Urbano",
			Description: pageSubtitle + " - Olho Urbano. Plataforma cidadã para reportar e acompanhar problemas urbanos. Transparência e participação social.",
			Keywords:    pageTitle + ", olho urbano, plataforma cidadã, denúncias urbanas, transparência pública, participação social, problemas urbanos, cidadania ativa",
			Canonical:   baseURL + "/" + pageSlug,
			Breadcrumbs: GenerateBreadcrumbs("footer_page", additionalData),
			PageSlug:    pageSlug,
		}
	default:
		return SEOData{
			Title:       "Olho Urbano - Plataforma de Denúncias Urbanas",
			Description: "Plataforma cidadã para reportar e acompanhar problemas urbanos.",
			Keywords:    "olho urbano, plataforma cidadã, denúncias urbanas",
			Canonical:   baseURL + "/",
			Breadcrumbs: GenerateBreadcrumbs("index", nil),
		}
	}
}
