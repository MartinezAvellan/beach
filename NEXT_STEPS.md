# NEXT STEPS - Security Fixes

## CRITICAL

- [ ] **Security Headers** - Adicionar middleware com `X-Content-Type-Options`, `X-Frame-Options`, `Strict-Transport-Security`, `Referrer-Policy` (`router.go`)
- [ ] **XSS via innerHTML** - Substituir `innerHTML` por `textContent` ou sanitizar com DOMPurify (`index.html`)
- [ ] **CORS/CSP** - Adicionar Content-Security-Policy restritivo e pinar HLS.js a versao exata (`router.go`, `index.html`)

## HIGH

- [ ] **Validacao do camera ID** - Validar tamanho e caracteres permitidos no path param `{id}` (`handlers.go`)
- [ ] **Rate limiting** - Adicionar middleware de rate limit por IP (`router.go`)
- [ ] **Limite de resposta IPMA** - Usar `io.LimitReader` no decode JSON da API externa (`client.go`)

## MEDIUM

- [ ] **Docker non-root** - Adicionar `USER` no Dockerfile (`Dockerfile`)
- [ ] **K8s securityContext** - Adicionar `runAsNonRoot`, `readOnlyRootFilesystem`, drop `ALL` capabilities (`deployment.yaml`)
- [ ] **MaxHeaderBytes** - Configurar limite no HTTP server (`main.go`)
- [ ] **Custom Recoverer** - Substituir `middleware.Recoverer` por handler que nao expoe stack traces (`router.go`)

## LOW

- [ ] **Pinar HLS.js** - Trocar `@1` por versao exata ex: `@1.4.12` (`index.html`)
- [ ] **Health check** - Validar conectividade com IPMA no `/health` (`handlers.go`)
- [ ] **TLS no Ingress** - Configurar TLS no Ingress do K8s (`deployment.yaml`)