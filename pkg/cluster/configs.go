package cluster

const odooDefaultSection = `
[DEFAULT]

without_demo = %s
server_wide_modules = %s

db_name = %s
db_template = %s

list_db = %s
dbfilter = %s

unaccent = True
publisher_warranty_url = %s

log_handler = %s

`
const odooOptionsSection = `
[options]
data_dir = %s
proxy_mode = True

# A static option used by custom in-app backup module
backupfolder = %s

# Cluster scoped overrides, if defined:
%s

# Track scoped overrides, if defined:
%s
`

const odooIntegratorSection = `

[integrator]

# Cluster scoped overrides, if defined:
%s

# Track scoped overrides, if defined:
%s

`

const odooCustomSection = `

# Cluster scoped custom sections, if defined:
%s

# Track scoped custom sections, if defined:
%s

`

const odooPsqlSecretFmt = `
# hostname:port:database:username:password
# If an entry needs to contain : or \, escape this character with \
%s:%s:*:%s:%s
`
const odooAdminSecretFmt = `%s

`
