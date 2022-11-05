#!/usr/bin/env zx --experimental

const expected = 90
const exclude = [
  'cmd',
  'checker/mock',
  'vm/func_types',
  'vm/runtime/helpers',
]

cd(path.resolve(__dirname, '..', '..'))

await spinner('Running tests', async () => {
  await $`go test -coverprofile=coverage.out -coverpkg=github.com/antonmedv/expr/... ./...`
  const coverage = fs.readFileSync('coverage.out').toString()
    .split('\n')
    .filter(line => {
      for (const ex of exclude)
        if (line.includes(ex)) return false
      return true
    })
    .join('\n')
  fs.writeFileSync('coverage.out', coverage)
  await $`go tool cover -html=coverage.out -o coverage.html`
})

const cover = await $`go tool cover -func=coverage.out`
const total = +cover.stdout.match(/total:\s+\(statements\)\s+(\d+\.\d+)%/)[1]
if (total < expected) {
  echo(chalk.red(`Coverage is too low: ${total}% < ${expected}% (expected)`))
  process.exit(1)
}
