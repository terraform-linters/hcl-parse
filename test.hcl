cond = var.enabled ? (true) : func(1, var.input)[0]

foo "var" "baz" {
  for = [for x in var.foo: x + 1 if x < 10]
  obj = { a = var.bar[*], var.foo = var.baz[var.qux], c = [1, 2] } 
  temp = "%{ for v in [true] }${v}%{ endfor }"
  wrap = "${true}"
}
