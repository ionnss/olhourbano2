# Olho Urbano

<div align="center">
  <img src="static/resource/og-image.png" alt="Olho Urbano Logo" width="500"/>
  <br><br>
  
  **Plataforma de Tecnologia Cívica Open Source**
  
  [![Licença: Transparência](https://img.shields.io/badge/Licença-Transparência-blue.svg)](LICENSE)
  [![Status: Produção](https://img.shields.io/badge/Status-Produção-brightgreen.svg)](https://olhourbano.com.br)
  [![Go Version](https://img.shields.io/badge/Go-1.24.4+-blue.svg)](go.mod)
  [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)](docker-compose.yml)
</div>

## Sobre

**Olho Urbano** é uma plataforma de tecnologia cívica que conecta cidadãos com seus governos locais para reportar e resolver problemas urbanos. Nossa missão é tornar as cidades mais transparentes, responsivas e amigáveis aos cidadãos através da tecnologia.

**Plataforma Ativa**: [https://olhourbano.com.br](https://olhourbano.com.br)

## Transparência e Segurança

Este repositório está **aberto para inspeção** para demonstrar nosso compromisso com:

- **Transparência de Segurança**: Todo o código está visível para auditoria de segurança
- **Proteção de Privacidade**: Práticas transparentes de manipulação de dados
- **Governo Aberto**: Transparência em tecnologia cívica
- **Construção de Confiança**: Verificação pública de nossas práticas

## Licença e Uso

### Usos Permitidos
- **Auditoria de Segurança**: Revisar código para vulnerabilidades e práticas de segurança
- **Verificação de Transparência**: Verificar nossas práticas de manipulação de dados e privacidade
- **Fins Educacionais**: Aprender com nossa implementação de tecnologia cívica
- **Relatórios de Bugs**: Contribuir com descobertas de segurança e relatórios de bugs
- **Desenvolvimento Local**: Testar e auditar localmente para fins de segurança

### Usos Restritos
- **Cópia/Redistribuição**: Este código não pode ser copiado ou redistribuído
- **Trabalhos Derivados**: Criar versões modificadas não é permitido
- **Uso Comercial**: Implantação comercial requer permissão explícita
- **Implantação Não Autorizada**: Executar este software requer consentimento por escrito

**Licença**: [Licença de Transparência](LICENSE) - Visualizar para transparência, restrito para cópia

## Arquitetura

### Stack Tecnológico
- **Backend**: Go 1.24.4+ com framework web customizado
- **Banco de Dados**: PostgreSQL 15+ com sistema de migração customizado
- **Frontend**: HTML5, CSS3, JavaScript (Vanilla)
- **Infraestrutura**: Docker, Docker Compose, Caddy 2.10.0 (proxy reverso)
- **Segurança**: Gerenciamento customizado de segredos, aplicação de HTTPS

### Principais Recursos
- **Verificação Segura de CPF**: Integração com validação oficial de CPF brasileiro
- **Mapas Interativos**: Integração com Google Maps para relatórios baseados em localização
- **Design Mobile-First**: Interface responsiva para todos os dispositivos
- **Estatísticas em Tempo Real**: Análises de relatórios e votação ao vivo
- **Processamento de Arquivos**: Suporte para imagens, vídeos, PDFs com limpeza de metadados
- **Recursos Comunitários**: Sistema de comentários e votação
- **Sistema de Artigos**: Blog integrado com suporte a Markdown

## Auditoria e Desenvolvimento Local

### Para Fins de Auditoria e Segurança

Este repositório pode ser usado para:

- **Auditoria de Segurança**: Revisar práticas de segurança e proteção de dados
- **Verificação de Transparência**: Confirmar nossas reivindicações sobre transparência
- **Pesquisa Acadêmica**: Estudar implementações de tecnologia cívica
- **Relatórios de Vulnerabilidades**: Identificar e reportar problemas de segurança
- **Testes de Penetração**: Realizar testes de segurança autorizados

### Pré-requisitos para Desenvolvimento Local

- Go 1.24.4+
- PostgreSQL 15+
- Docker & Docker Compose
- ImageMagick (para processamento de arquivos)

### Configuração para Auditoria Local

#### 1. Clone o Repositório
```bash
git clone <repository-url>
cd olhourbano2
```

#### 2. Configuração de Segredos (Para Testes)
Crie os arquivos de segredo necessários em `secrets/`:
```bash
mkdir -p secrets
echo "test_db_password" > secrets/db_password.txt
echo "test_smtp_password" > secrets/smtp_password.txt
echo "test_session_key" > secrets/session_key.txt
echo "test_cpfhub_api_key" > secrets/cpfhub_api_key.txt
echo "test_google_maps_api_key" > secrets/google_maps_api_key.txt
```

#### 3. Variáveis de Ambiente (Para Testes)
Configure as variáveis de ambiente para desenvolvimento local:
```bash
export POSTGRES_DB=test_database
export POSTGRES_USER=test_user
export APP_VERSION=2.0.0
export SMTP_HOST=smtp.gmail.com
export SMTP_PORT=587
export SMTP_USERNAME=test_email@gmail.com
export COOKIE_DOMAIN=localhost
```

> **⚠️ Segurança**: Use apenas dados de teste para desenvolvimento local. Nunca use credenciais reais em ambiente de auditoria.

#### 4. Executar com Docker (Para Auditoria)
```bash
# Construir e iniciar todos os serviços
docker compose up -d

# Verificar status dos serviços
docker compose ps

# Visualizar logs
docker compose logs -f backend
```

#### 5. Migrações do Banco de Dados
```bash
# Executar migrações
docker exec -w /app your-backend-container /usr/local/bin/app migrate

# Verificar status das migrações
docker exec -w /app your-backend-container /usr/local/bin/app migrate:status

# Validar migrações
docker exec -w /app your-backend-container /usr/local/bin/app migrate:validate
```

#### 6. Acessar a Aplicação Local
- **Local**: http://localhost:8081
- **Produção**: https://olhourbano.com.br

## Estrutura do Projeto

```
olhourbano2/
├── articles/                 # Artigos em Markdown
├── config/                   # Configurações e categorias
├── db/                       # Banco de dados e migrações
├── handlers/                 # Handlers HTTP e endpoints da API
├── middleware/               # Middlewares customizados
├── models/                   # Estruturas de dados
├── routes/                   # Definição de rotas
├── secrets/                  # Arquivos de segredo
├── services/                 # Lógica de negócio
├── static/                   # Arquivos estáticos (CSS, JS, imagens)
├── templates/                # Templates HTML
├── uploads/                  # Arquivos enviados pelos usuários
├── docker-compose.yml        # Configuração Docker
├── Dockerfile               # Imagem Docker
├── Caddyfile                # Configuração do proxy reverso
└── main.go                  # Ponto de entrada da aplicação
```

## Componentes Principais

### Módulos Principais
- **`handlers/`**: Handlers de requisições HTTP e endpoints da API
- **`services/`**: Lógica de negócio e integrações com serviços externos
- **`models/`**: Estruturas de dados e modelos do banco de dados
- **`config/`**: Gerenciamento de configuração e categorias
- **`db/`**: Conexão com banco de dados e sistema de migração
- **`templates/`**: Templates HTML e componentes da interface

### Recursos de Segurança
- **Gerenciamento de Segredos**: Segredos baseados em arquivo com integração Docker
- **Verificação de CPF**: API oficial de validação de CPF brasileiro
- **Segurança de Arquivos**: Limpeza de metadados e validação de tipos
- **Aplicação de HTTPS**: SSL/TLS com headers de segurança
- **Validação de Entrada**: Validação abrangente de formulários

## Documentação

- **[Política de Transparência](TRANSPARENCY.md)**: Nossa abordagem à transparência e código aberto
- **[Guia de Contribuição](CONTRIBUTING.md)**: Como contribuir com o projeto
- **[Guia de Migração](MIGRATION_GUIDE.md)**: Migração entre servidores e backup
- **[Esquema do Banco de Dados](db/README.md)**: Estrutura do banco de dados e sistema de migração
- **[Configuração](config/)**: Arquivos de configuração e categorias

## Impacto da Tecnologia Cívica

O Olho Urbano demonstra como a tecnologia cívica pode:

- **Empoderar Cidadãos**: Fornecer canais fáceis para relatar problemas urbanos
- **Melhorar a Transparência**: Visibilidade pública das respostas governamentais
- **Aumentar a Responsabilidade**: Acompanhar resolução de problemas e tempos de resposta
- **Fomentar a Comunidade**: Permitir colaboração cidadã e votação

## Contribuindo

### Segurança e Relatórios de Bugs
Acolhemos descobertas de segurança e relatórios de bugs:
- **Problemas de Segurança**: Reporte via email para olhourbano.contato@gmail.com
- **Relatórios de Bugs**: Use GitHub Issues para bugs não relacionados à segurança
- **Documentação**: Sugira melhorias em nossa documentação

### Contribuições de Código
Devido ao nosso modelo de licenciamento, não podemos aceitar contribuições diretas de código. No entanto, acolhemos:
- **Auditorias de Segurança**: Revisões independentes de segurança
- **Sugestões de Recursos**: Ideias para melhorias da plataforma
- **Documentação**: Ajude a melhorar nossa documentação de transparência

## Contato e Suporte

- **Email Principal**: olhourbano.contato@gmail.com
- **Email Secundário**: zara.leonardo@gmail.com
- **Website**: [https://olhourbano.com.br](https://olhourbano.com.br)

## Legal

Este software é fornecido sob a [Licença de Transparência](LICENSE). 
Para consultas sobre licenciamento ou solicitações de uso comercial, entre em contato: olhourbano.contato@gmail.com

---

**Olho Urbano** - Construindo cidades transparentes e responsivas através da tecnologia cívica.

*"A transparência constrói confiança. A confiança constrói cidades melhores."*