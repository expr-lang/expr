#!/usr/bin/env zx

const expected = 90
const exclude = [
  'expr/test', // We do not need to test the test package.
  'checker/mock', // Mocks only used for testing.
  'vm/func_types', // Generated files.
  'vm/runtime/helpers', // Generated files.
  'internal/difflib', // Test dependency. This is vendored dependency, and ideally we also have good tests for it.
  'internal/spew', // Test dependency.
  'internal/testify', // Test dependency.
  'patcher/value', // Contains a lot of repeating code. Ideally we should have a test for it.
  'pro', // Expr Pro is not a part of the main codebase.
]

cd(path.resolve(__dirname, '..', '..'))

await spinner('Running tests', async () => {
  await $`go test -coverprofile=coverage.out -coverpkg=github.com/expr-lang/expr/... ./...`
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

const cover = await $({verbose: true})`go tool cover -func=coverage.out`
const total = +cover.stdout.match(/total:\s+\(statements\)\s+(\d+\.\d+)%/)[1]
if (total < expected) {
  echo(chalk.red(`Coverage is too low: ${total}% < ${expected}% (expected)`))
  process.exit(1)
} else {
  echo(`Coverage is good: ${chalk.green(total + '%')} >= ${expected}% (expected)`)
}
