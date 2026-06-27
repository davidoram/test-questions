# Dependency Resolver Fixtures

Each package is represented by a directory containing a simplified `package.json`.

Suggested starting packages:

- `app-basic`: simple acyclic dependency chain.
- `app-shared`: acyclic graph with a shared transitive dependency.
- `app-dev`: includes `devDependencies` so you can choose whether your resolver includes them.
- `app-cycle`: intentionally depends on a cycle and should fail cycle detection.

One valid install order for `app-basic` is:

```text
delta
charlie
beta
alpha
app-basic
```

One valid install order for `app-shared` is:

```text
delta
charlie
echo
alpha
app-shared
```

For `app-cycle`, your resolver should detect:

```text
cycle-a -> cycle-b -> cycle-c -> cycle-a
```
