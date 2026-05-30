setup() {
    load 'test_helper/bats-support/load'
    load 'test_helper/bats-assert/load'
    load 'test_helper/bats-file/load'

    TMP_DIR="/tmp/framed-imports-test"
    DIR="$( cd "$( dirname "$BATS_TEST_FILENAME" )" >/dev/null 2>&1 && pwd )"

    mkdir -p "$TMP_DIR"
    cd "$TMP_DIR"
}

teardown() {
    sudo rm -R "$TMP_DIR"
}

@test "Local template import is resolved at load time" {
    cat > shared.yaml <<'YAML'
name: shared
structure:
  name: ignored
  files:
    - common.md
  dirs:
    - name: nested
      files:
        - inner.md
YAML

    cat > framed.yaml <<'YAML'
name: parent
structure:
  name: root
  dirs:
    - name: lib
      template: ./shared.yaml
YAML

    run framed visualize
    assert_success
    assert_output --partial 'lib'
    assert_output --partial 'nested'
    assert_output --partial 'common.md'
    assert_output --partial 'inner.md'
}

@test "Verify works against a template assembled from a local import" {
    cat > shared.yaml <<'YAML'
name: shared
structure:
  name: ignored
  files:
    - common.md
YAML

    cat > framed.yaml <<'YAML'
name: parent
structure:
  name: root
  dirs:
    - name: lib
      template: ./shared.yaml
YAML

    mkdir -p lib
    touch lib/common.md

    run framed verify
    assert_success
    assert_output --partial '✅ Verified successfully!'
}

@test "Missing imported file surfaces a clear error" {
    cat > framed.yaml <<'YAML'
name: parent
structure:
  name: root
  dirs:
    - name: oops
      template: ./does-not-exist.yaml
YAML

    run framed verify
    assert_failure
}

@test "Declaring both template and inline files is rejected" {
    cat > shared.yaml <<'YAML'
name: shared
structure:
  name: x
  files:
    - hello.md
YAML

    cat > framed.yaml <<'YAML'
name: parent
structure:
  name: root
  dirs:
    - name: hybrid
      template: ./shared.yaml
      files:
        - local.md
YAML

    run framed verify
    assert_failure
}

@test "Cyclic local imports are detected" {
    cat > a.yaml <<'YAML'
name: a
structure:
  name: x
  dirs:
    - name: hop
      template: ./b.yaml
YAML

    cat > b.yaml <<'YAML'
name: b
structure:
  name: x
  dirs:
    - name: back
      template: ./a.yaml
YAML

    cat > framed.yaml <<'YAML'
name: parent
structure:
  name: root
  dirs:
    - name: start
      template: ./a.yaml
YAML

    run framed verify
    assert_failure
}
