# Templates Directory Structure

This document explains the template organization for the Olho Urbano project and how to use Go's `html/template` package for creating reusable HTML components.

## Directory Structure

```
templates/
├── README.md           # This documentation file
├── layouts/           # Base HTML layouts (main page structure)
├── components/        # Reusable HTML components
└── pages/            # Page-specific content templates
```

## Template Types Explained

### 1. Layouts (`layouts/`)

**Purpose**: Define the main HTML structure that wraps all pages.

**Contains**:
- Complete HTML document structure (`<!DOCTYPE>`, `<html>`, `<head>`, `<body>`)
- Meta tags, page title, CSS/JS links
- Common page structure elements
- Template blocks where content gets inserted

**When to use**: Create a layout when you need a different overall page structure (e.g., admin layout vs public layout, login page layout).

**Example file**: `layouts/base.html`
```html
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - Olho Urbano</title>
    <link rel="stylesheet" href="/static/css/main.css">
</head>
<body>
    {{template "header" .}}
    
    <main class="main-content">
        {{block "content" .}}{{end}}
    </main>
    
    {{template "footer" .}}
    
    <script src="/static/js/app.js"></script>
</body>
</html>
```

### 2. Components (`components/`)

**Purpose**: Small, reusable HTML pieces that can be included in any layout or page.

**Contains**:
- Headers, footers, navigation bars
- Sidebars, modals, forms, alerts
- Any HTML snippet used in multiple places
- Each component is wrapped in `{{define "componentName"}}`

**When to use**: Create a component when you have HTML that will be used in multiple pages or layouts.

**Example files**:

**`components/header.html`**:
```html
{{define "header"}}
<header class="main-header">
    <div class="container">
        <div class="header-brand">
            <img src="/static/resource/full_logo.svg" alt="Olho Urbano" class="logo">
            <h1>Olho Urbano</h1>
        </div>
        <nav class="main-nav">
            <a href="/" {{if eq .CurrentPage "home"}}class="active"{{end}}>Home</a>
            <a href="/reports" {{if eq .CurrentPage "reports"}}class="active"{{end}}>Relatórios</a>
            <a href="/dashboard" {{if eq .CurrentPage "dashboard"}}class="active"{{end}}>Dashboard</a>
        </nav>
        <div class="user-info">
            {{if .User}}
                <span>Olá, {{.User.Name}}</span>
                <a href="/logout">Sair</a>
            {{else}}
                <a href="/login">Entrar</a>
            {{end}}
        </div>
    </div>
</header>
{{end}}
```

**`components/footer.html`**:
```html
{{define "footer"}}
<footer class="main-footer">
    <div class="container">
        <div class="footer-content">
            <div class="footer-section">
                <h3>Olho Urbano</h3>
                <p>Transparência e participação cidadã para um futuro melhor.</p>
            </div>
            <div class="footer-section">
                <h4>Links Úteis</h4>
                <ul>
                    <li><a href="/about">Sobre</a></li>
                    <li><a href="/privacy">Privacidade</a></li>
                    <li><a href="/terms">Termos de Uso</a></li>
                </ul>
            </div>
            <div class="footer-section">
                <h4>Contato</h4>
                <p>Email: contato@olhourbano.com.br</p>
            </div>
        </div>
        <div class="footer-bottom">
            <p>&copy; {{.CurrentYear}} Olho Urbano. Todos os direitos reservados.</p>
        </div>
    </div>
</footer>
{{end}}
```

**`components/navigation.html`**:
```html
{{define "navigation"}}
<nav class="sidebar-nav">
    <ul>
        <li><a href="/dashboard" {{if eq .CurrentPage "dashboard"}}class="active"{{end}}>
            <i class="icon-dashboard"></i> Dashboard
        </a></li>
        <li><a href="/reports" {{if eq .CurrentPage "reports"}}class="active"{{end}}>
            <i class="icon-reports"></i> Relatórios
        </a></li>
        <li><a href="/complaints" {{if eq .CurrentPage "complaints"}}class="active"{{end}}>
            <i class="icon-complaints"></i> Denúncias
        </a></li>
        <li><a href="/analytics" {{if eq .CurrentPage "analytics"}}class="active"{{end}}>
            <i class="icon-analytics"></i> Análises
        </a></li>
    </ul>
</nav>
{{end}}
```

### 3. Pages (`pages/`)

**Purpose**: Contain the unique content for each specific page/route.

**Contains**:
- Page-specific HTML content
- Forms, data displays, content sections
- Content that fills the layout's "content" block
- Each page defines a `{{define "content"}}` block

**When to use**: Create a page template for each unique route/URL in your application.

**Example files**:

**`pages/home.html`**:
```html
{{define "content"}}
<div class="hero-section">
    <div class="hero-content">
        <h1>Bem-vindo ao Olho Urbano</h1>
        <p class="hero-subtitle">Transparência e participação cidadã para um futuro melhor</p>
        <div class="hero-actions">
            <a href="/reports" class="btn btn-primary">Ver Relatórios</a>
            <a href="/about" class="btn btn-secondary">Saiba Mais</a>
        </div>
    </div>
</div>

<div class="features-section">
    <div class="container">
        <h2>O que oferecemos</h2>
        <div class="features-grid">
            <div class="feature-card">
                <img src="/static/resource/blue_eye.svg" alt="Relatórios">
                <h3>Relatórios Públicos</h3>
                <p>Acesse informações detalhadas sobre gastos e projetos públicos da sua cidade.</p>
            </div>
            <div class="feature-card">
                <img src="/static/resource/grey_eye.svg" alt="Denúncias">
                <h3>Canal de Denúncias</h3>
                <p>Reporte problemas urbanos e acompanhe o progresso das soluções.</p>
            </div>
            <div class="feature-card">
                <img src="/static/resource/circular_eye.png" alt="Transparência">
                <h3>Transparência Total</h3>
                <p>Dados abertos e acessíveis para todos os cidadãos.</p>
            </div>
        </div>
    </div>
</div>
{{end}}
```

**`pages/dashboard.html`**:
```html
{{define "content"}}
<div class="dashboard-header">
    <h1>Dashboard</h1>
    <p>Visão geral dos dados e atividades</p>
</div>

<div class="dashboard-content">
    <div class="stats-grid">
        <div class="stat-card">
            <h3>{{.Stats.TotalReports}}</h3>
            <p>Relatórios Publicados</p>
        </div>
        <div class="stat-card">
            <h3>{{.Stats.PendingComplaints}}</h3>
            <p>Denúncias Pendentes</p>
        </div>
        <div class="stat-card">
            <h3>{{.Stats.ActiveUsers}}</h3>
            <p>Usuários Ativos</p>
        </div>
    </div>

    <div class="dashboard-sections">
        <section class="recent-reports">
            <h2>Relatórios Recentes</h2>
            {{range .RecentReports}}
            <div class="report-item">
                <h4>{{.Title}}</h4>
                <p>{{.Summary}}</p>
                <span class="date">{{.CreatedAt.Format "02/01/2006"}}</span>
            </div>
            {{end}}
        </section>

        <section class="quick-actions">
            <h2>Ações Rápidas</h2>
            <div class="action-buttons">
                <a href="/reports/new" class="btn btn-primary">Novo Relatório</a>
                <a href="/complaints/new" class="btn btn-secondary">Nova Denúncia</a>
                <a href="/analytics" class="btn btn-info">Ver Análises</a>
            </div>
        </section>
    </div>
</div>
{{end}}
```

## How Templates Work Together

### Template Execution Flow

1. **Handler receives request** for `/dashboard`
2. **Load templates**: Layout + Components + Page
3. **Parse data**: Prepare data to pass to template
4. **Execute template**: Combine all parts with data
5. **Send response**: Complete HTML to browser

### Template Composition Example

```
base.html (layout)
├── header.html (component)     ← {{template "header" .}}
├── dashboard.html (page)       ← {{block "content" .}}
└── footer.html (component)     ← {{template "footer" .}}
```

### Data Passing

Templates can access data passed from Go handlers:

```go
type PageData struct {
    Title       string
    CurrentPage string
    User        *User
    Stats       *DashboardStats
    CurrentYear int
}

data := PageData{
    Title:       "Dashboard",
    CurrentPage: "dashboard",
    User:        getCurrentUser(),
    Stats:       getDashboardStats(),
    CurrentYear: time.Now().Year(),
}

tmpl.ExecuteTemplate(w, "dashboard.html", data)
```

## Go Handler Implementation

### Basic Template Loading

```go
// In your main.go or template package
var templates *template.Template

func init() {
    // Load all templates
    templates = template.Must(template.ParseGlob("templates/**/*.html"))
}

// In your handlers
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
    data := PageData{
        Title:       "Dashboard",
        CurrentPage: "dashboard",
        User:        getCurrentUser(r),
        Stats:       getDashboardStats(),
    }
    
    err := templates.ExecuteTemplate(w, "base.html", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
```

## Best Practices

### 1. Naming Conventions
- **Layouts**: `base.html`, `admin.html`, `login.html`
- **Components**: `header.html`, `footer.html`, `nav.html`
- **Pages**: Match your route names (`home.html`, `dashboard.html`, `reports.html`)

### 2. Template Definitions
- **Components**: Always wrap in `{{define "componentName"}}`
- **Pages**: Always define `{{define "content"}}`
- **Layouts**: Use `{{block "content" .}}` for page insertion

### 3. Data Structure
- Create consistent data structures for common elements
- Pass `CurrentPage` to enable active navigation states
- Include user info for authentication-dependent components

### 4. Performance
- Parse templates once at startup, not on every request
- Use template caching in production
- Minimize template complexity for better performance

## Common Template Functions

Go templates provide useful built-in functions:

```html
<!-- Conditionals -->
{{if .User}}Welcome, {{.User.Name}}{{end}}

<!-- Loops -->
{{range .Items}}
    <li>{{.Name}}</li>
{{end}}

<!-- Comparisons -->
{{if eq .CurrentPage "home"}}class="active"{{end}}

<!-- Date formatting -->
{{.CreatedAt.Format "02/01/2006 15:04"}}

<!-- HTML safety (use with caution) -->
{{.HTMLContent | html}}
```

## Security Notes

- **Never use `| html` unless absolutely necessary** - Go templates auto-escape by default
- **Validate all user input** before passing to templates
- **Don't expose sensitive data** in template variables
- **Use HTTPS** for any pages with user data

---

*This documentation covers the template structure for Olho Urbano. For implementation examples, see the handlers/ directory.* 