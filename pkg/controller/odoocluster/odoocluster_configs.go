/*
 * This file is part of the Odoo-Operator (R) project.
 * Copyright (c) 2018-2018 XOE Corp. SAS
 * Authors: David Arnold, et al.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *
 * ALTERNATIVE LICENCING OPTION
 *
 * You can be released from the requirements of the license by purchasing
 * a commercial license. Buying such a license is mandatory as soon as you
 * develop commercial activities involving the Odoo-Operator software without
 * disclosing the source code of your own applications. These activities
 * include: Offering paid services to a customer as an ASP, shipping Odoo-
 * Operator with a closed source product.
 *
 */

package odoocluster

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
