# Examples

Character Frequency Grouping

```
let char = "0";
1..100000
| map(string(#))
| groupBy(split(#, char) | len() - 1)
| toPairs()
| map({{
    count: #[0], 
    len: len(#[1]), 
    examples: [first(#[1]), last(#[1])],
 }})
| sortBy(.len, 'desc')
```

Log Filtering and Aggregation

```
let logs = [
  {timestamp: date("2023-08-14 08:30:00"), message: "User logged in", level: "info"},
  {timestamp: date("2023-08-14 09:00:00"), message: "Error processing payment", level: "error"},
  {timestamp: date("2023-08-14 10:15:00"), message: "User logged out", level: "info"},
  {timestamp: date("2023-08-14 11:00:00"), message: "Error connecting to database", level: "error"}
];

logs
| filter(.level == "error")
| map({{
    time: string(.timestamp),
    detail: .message
 }})
| sortBy(.time)
```

Financial Data Analysis and Summary

```
let accounts = [
  {name: "Alice", balance: 1234.56, transactions: [100, -50, 200]},
  {name: "Bob", balance: 2345.67, transactions: [-200, 300, -150]},
  {name: "Charlie", balance: 3456.78, transactions: [400, -100, 50]}
];

{
  totalBalance: sum(accounts, .balance),
  averageBalance: mean(map(accounts, .balance)),
  totalTransactions: reduce(accounts, #acc + len(.transactions), 0),
  accounts: map(accounts, {{
      name: .name,
      final: .balance + sum(.transactions),
      transactionCount: len(.transactions)
  }})
}
```

Bitwise Operations and Flags Decoding

```
let flags = [
  {name: "read", value: 0b0001},
  {name: "write", value: 0b0010},
  {name: "execute", value: 0b0100},
  {name: "admin", value: 0b1000}
];

let userPermissions = 0b1011;

flags
| filter(userPermissions | bitand(.value) != 0)
| map(.name)
```

Nested Predicates with Optional Chaining

```
let users = [
  {id: 1, name: "Alice", posts: [{title: "Hello World", content: "Short post"}, {title: "Another Post", content: "This is a bit longer post"}]},
  {id: 2, name: "Bob", posts: nil},
  {id: 3, name: "Charlie", posts: [{title: "Quick Update", content: "Update content"}]}
];

users
| filter(
    // Check if any post has content length greater than 10.
    any(.posts ?? [], len(.content) > 10)
  )
| map({{name: .name, postCount: len(.posts ?? [])}})
```

String Manipulation and Validation

```
"  Apple, banana, Apple, orange, banana, kiwi "
| trim()
| split(",")
| map(trim(#))
| map(lower(#))
| uniq()
| sort()
| join(", ")
```

Date Difference

```
let startDate = date("2023-01-01");
let endDate = date("2023-12-31");
let diff = endDate - startDate;
{
  daysBetween: diff.Hours() / 24,
  hoursBetween: diff.Hours()
}
```

Phone number filtering

```
let phone = filter(split("123-456-78901", ""), # in map(0..9, string(#)))[:10];
join(concat(["("], phone[:3], [")"], phone[3:6], ["-"], phone[6:]))
```

Prime numbers

```
2..1000 | filter(let N = #; none(2..N-1, N % # == 0))
```
