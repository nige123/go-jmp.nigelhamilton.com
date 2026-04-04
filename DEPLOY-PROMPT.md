# Deploy Setup Prompt

Use this prompt with Claude Code in any Mojolicious repo under /home/s3/ to set up standardised deployment infrastructure.

---

## Prompt

```
Set up deployment infrastructure for this Mojolicious repo. Follow these conventions exactly:

### 1. Nginx config

Create `conf/nginx/{domain}.conf` (production) with this pattern:

- upstream block pointing to the hypnotoad port from conf/production.conf
- HTTP (80) → HTTPS (301) redirect
- HTTPS (443) with SSL from /etc/letsencrypt/live/{domain}/
- TLSv1.2 + TLSv1.3, HIGH ciphers
- proxy_pass to upstream with headers: Host, X-Real-IP, X-Forwarded-For, X-Forwarded-Proto, Upgrade, Connection
- access_log and error_log to /var/log/nginx/{domain}.*.log

If a test config already exists in conf/nginx/, keep it. Model production config on /home/s3/api.wh.ax/conf/nginx.conf.

Print the symlink command: sudo ln -sf $(pwd)/conf/nginx/{domain}.conf /etc/nginx/sites-enabled/{domain}

### 2. Ubic service

Ensure ubic service files exist in ubic/service/{category}/{service_name} using Ubic::Service::SimpleDaemon with:

- cwd set to the repo root
- bin using: perlbrew exec --with perl-5.42.1 hypnotoad -f {repo}/bin/app.pl
- stdout/stderr/ubic_log to /tmp/{domain}.{service}.{stdout|stderr|ubic}.log

Print the symlink command: ln -sf $(pwd)/ubic/service/{category} ~/ubic/service/{category}

### 3. Admin deploy route

Add a password-protected `/admin/deploy` route (NOT a subdomain action — a plain path route).

a) Add a new route in the app startup (or RoutesConfig):
   - GET  /admin/deploy → Controller::Admin#show_deploy
   - POST /admin/deploy → Controller::Admin#do_deploy

b) Create lib/{App}/Controller/Admin.pm (a standard Mojolicious controller, NOT an action handler) with these methods:

   - show_deploy (GET): renders a deploy status page showing current git SHA, branch, last deploy time. Include a "Deploy Now" button (POST form). Password-protect with a simple token check — read deploy_token from config, compare against a "token" param stored in the session. Show a token input form if not authenticated.

   - do_deploy (POST): checks deploy_token first, then executes deployment steps in order, capturing output:
     1. git pull origin {main_branch}
     2. cpanm --local-lib local --notest --installdeps .
     3. ubic restart {service_name}
     4. Record deploy timestamp
     Render results page showing success/failure and command output for each step.

c) Add to config files:
   - deploy_token: a secret token (use $ENV{DEPLOY_TOKEN} with a fallback)
   - deploy_service: the ubic service name to restart
   - deploy_branch: the git branch to pull (default: main or master)

d) Add the /admin/deploy routes BEFORE the catch-all wildcard route so they take priority.

e) Create the EP templates for the deploy pages:
   - pages/admin-deploy.html.ep — token form + status + deploy button (use tok() helper, no boilerplate)
   - pages/admin-deploy-result.html.ep — shows step-by-step output

f) Security: the deploy route must check the deploy_token before allowing any action. Use session-based auth so the token is entered once per browser session. Never execute shell commands with user-supplied input — all commands are hardcoded paths.

### 4. Update install.pl

Add to the existing install.pl:
- Print the nginx symlink command
- Print the ubic symlink command
- Don't auto-run sudo commands, just print them for the user

### Conventions to follow

- Perl with Mojo::Base -signatures
- use strict; use warnings; everywhere
- Comment blocks above each method with separator lines
- Role::Tiny for handler roles
- Templates use tok() and tok_raw() Mojo helpers (no boilerplate closures)
- One API call per page max
- Config values via $self->config->{key}
- All shell commands in deploy handler use absolute paths
- Log output with $c->app->log->info()

### Do NOT

- Add any JavaScript frameworks
- Make client-side API calls
- Add external dependencies beyond what's in cpanfile
- Hardcode passwords in source (use config with env var fallback)
- Run deploy commands with any user-supplied input in the shell command
```

---

## After running the prompt

1. Create the nginx symlink: `sudo ln -sf $(pwd)/conf/nginx/{domain}.conf /etc/nginx/sites-enabled/{domain}`
2. Reload nginx: `sudo nginx -t && sudo systemctl reload nginx`
3. Create the ubic symlink: `ln -sf $(pwd)/ubic/service/{category} ~/ubic/service/{category}`
4. Add DEPLOY_TOKEN to your environment or config
5. Test: visit `https://{domain}/admin/deploy` and enter your token
