#!/usr/bin/env zx --experimental

cd(path.resolve(__dirname, '..'))

await spinner('Running tests', async () => {
  await $`go test -coverprofile=coverage.out -coverpkg=github.com/antonmedv/expr/... ./...`
  const coverage = fs.readFileSync('coverage.out').toString()
    .split('\n')
    .filter(line => !line.match(/cmd|generate/))
    .join('\n')
  fs.writeFileSync('coverage.out', coverage)
  await $`go tool cover -html=coverage.out -o coverage.html`
})

await $`go tool cover -func=coverage.out`
