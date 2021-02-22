# Registry-check

A tool for inspecting which registries are being used by package
managers (`npm` and `composer`).

# Do I need it?

If you use a private registry for installing packages you might want to have
an automated check in CI to ensure no packages are installed directly from
public registries. This tool will make it a little bit easier to do.

# Usage

```
registry-check provides a list of registry URLs used in a given NPM or Composer lockfile in text or JSON format.

USAGE:
  registry-check [OPTIONS] LOCKFILE

ARGS:
  <LOCKFILE>   package-lock.json or composer.lock file.

OPTIONS:
  -json
      Enable JSON output. (defaults to false)
  -type string
      Force lockfile type. Possible values: "npm", "composer". (defaults to guessing from filename)

EXAMPLES:
  registry-check composer.lock                    List registries in composer.lock. Output in text format (default).

  registry-check -json package-lock.json          List registries in package-lock.json. Outputs JSON.

  registry-check -type npm -json mylock.json      List registries in package-lock.json and force NPM lockfile format
```

# Simple example

## NPM project

### Plain text output

```
$ registry-check package-lock.json

Found 2 registries in package-lock.json with 743 packages:
- https://registry.npmjs.org
- http://registry.npmjs.org
```

### JSON output

JSON output is condensed, jq used for clarity.

```
$ registry-check -json package-lock.json | jq -r

{
  "success": true,
  "type": "npm",
  "packages": 743,
  "filename": "package-lock.json",
  "registries": [
    "https://registry.npmjs.org",
    "http://registry.npmjs.org"
  ]
}
```

## Combined examples

### Find all packages with more than one registry

#### Using JSON output with fd, xargs and jq.

```
$ fd package-lock.json -t f -0 \
  | xargs -0 -n1 registry-check -json \
  | jq -r '. | select(.registries | length > 1) | .filename, .registries'

package-a/package-lock.json
[
  "https://registry.npmjs.org",
  "http://registry.npmjs.org"
]
package-b/package-lock.json
[
  "https://registry.example.com",
  "https://registry.npmjs.org"
]
```

#### Using text output with fd, xargs and awk.

```
$ fd package-lock.json -t f -0 \
  | xargs -0 -n1 registry-check \
  | awk '$1=="Found" { p=0 }; $1=="Found" && $2 > 1 { p=1; print $5 }; $1=="-" && p { print }'

package-a/package-lock.json
- https://registry.npmjs.org
- http://registry.npmjs.org
package-b/package-lock.json
- https://registry.example.com
- https://registry.npmjs.org
```

### Find all packages NOT using your registry

#### Using JSON output with fd, xargs and jq.

```
$ fd package-lock.json -t f -0 \
  | xargs -0 -n1 registry-check -json \
  | jq -r '. | select(.registries | map(select(. != "https://registry.example.com")) | length > 0)'

package-a/package-lock.json
[
  "https://registry.npmjs.org",
  "http://registry.npmjs.org"
]
```

#### Using text output with fd, xargs, awk and uniq.

```
$ fd package-lock.json -t f -0  \
  | xargs -0 -n1 registry-check \
  | awk '$1=="Found" { filename=$5 }; $1 == "-" && $2 != "https://registry.example.com" { print filename }' \
  | uniq

package-a/package-lock.json
```
