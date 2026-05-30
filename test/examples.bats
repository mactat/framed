setup() {
    load 'test_helper/bats-support/load'
    load 'test_helper/bats-assert/load'
    load 'test_helper/bats-file/load'

    TMP_DIR="/tmp/framed-examples-test"
    REPO_ROOT="$( cd "$( dirname "$BATS_TEST_FILENAME" )/.." >/dev/null 2>&1 && pwd )"

    mkdir -p "$TMP_DIR"
    cd "$TMP_DIR"
}

teardown() {
    sudo rm -R "$TMP_DIR"
}

# Iterates through every example yaml the repo ships and confirms a
# fresh `framed create --files` followed by `framed verify` succeeds.
# Drives the same round-trip the README promises.
@test "Every shipped example survives create + verify" {
    for example in "$REPO_ROOT"/examples/*.yaml; do
        name="$(basename "$example" .yaml)"
        workdir="$TMP_DIR/$name"
        mkdir -p "$workdir"
        cp "$example" "$workdir/framed.yaml"

        ( cd "$workdir" && framed create --files >/dev/null )
        run bash -c "cd '$workdir' && framed verify"
        assert_success "example $name failed verify: $output"
        assert_output --partial '✅ Verified successfully!'
    done
}

@test "node example demonstrates allowed and forbidden patterns" {
    cp "$REPO_ROOT/examples/node.yaml" framed.yaml
    framed create --files >/dev/null
    touch src/sneaky.test.ts

    run framed verify
    assert_failure
    assert_output --partial 'Forbidden pattern'
}
