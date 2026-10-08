package grit

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// setupGritEnv creates a temp directory, changes into it, and initializes a
// complete grit environment (.grit dir, config.yml, history.log).
func setupGritEnv(t *testing.T) {
	t.Helper()
	t.Chdir(t.TempDir())

	if err := os.Mkdir(GritDir, 0755); err != nil {
		t.Fatal(err)
	}
	cfg, err := DefaultConfig()
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteConfig(cfg); err != nil {
		t.Fatal(err)
	}

	f, err := os.Create(HistoryFile)
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
}

// makeFakeGitDir creates a directory with a .git child (no real git init needed
// for tests that only check for .git existence).
func makeFakeGitDir(t *testing.T, name string) {
	t.Helper()
	if err := os.Mkdir(name, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(name, ".git"), 0755); err != nil {
		t.Fatal(err)
	}
}

// makeRealGitRepo runs git init to produce a repository that accepts git commands.
func makeRealGitRepo(t *testing.T, name string) {
	t.Helper()
	if err := os.Mkdir(name, 0755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("git", "init", name)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init %s: %v\n%s", name, err, out)
	}
}

// --- Header / Footer (pure functions) ---

func TestHeader_NoArgs(t *testing.T) {
	got := Header()
	want := "----------------------------------------\n"
	if got != want {
		t.Errorf("Header() = %q, want %q", got, want)
	}
}

func TestHeader_WithArg(t *testing.T) {
	got := Header("my-repo")
	want := "----------------------------------------\n>> my-repo"
	if got != want {
		t.Errorf("Header(%q) = %q, want %q", "my-repo", got, want)
	}
}

func TestFooter_NoArgs(t *testing.T) {
	got := Footer()
	want := "----------------------------------------\n"
	if got != want {
		t.Errorf("Footer() = %q, want %q", got, want)
	}
}

func TestFooter_WithArg(t *testing.T) {
	got := Footer("my-repo")
	want := "<< my-repo\n----------------------------------------\n"
	if got != want {
		t.Errorf("Footer(%q) = %q, want %q", "my-repo", got, want)
	}
}

// --- DefaultConfig ---

func TestDefaultConfig(t *testing.T) {
	t.Chdir(t.TempDir())

	config, err := DefaultConfig()
	if err != nil {
		t.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if config.Root != cwd {
		t.Errorf("Root = %q, want %q", config.Root, cwd)
	}
	if !config.IgnoreRoot {
		t.Error("IgnoreRoot should default to true")
	}
	if len(config.Repositories) != 0 {
		t.Errorf("Repositories should be empty, got %d entries", len(config.Repositories))
	}
}

// --- WriteConfig / LoadConfig round-trip ---

func TestWriteLoadConfigRoundTrip(t *testing.T) {
	setupGritEnv(t)

	want := Config{
		Root:       "/some/root",
		IgnoreRoot: false,
		Repositories: []Repository{
			{Name: "alpha", Path: "alpha"},
			{Name: "beta", Path: "beta"},
		},
	}
	if err := WriteConfig(want); err != nil {
		t.Fatal(err)
	}
	got, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	if got.Root != want.Root {
		t.Errorf("Root: got %q, want %q", got.Root, want.Root)
	}
	if got.IgnoreRoot != want.IgnoreRoot {
		t.Errorf("IgnoreRoot: got %v, want %v", got.IgnoreRoot, want.IgnoreRoot)
	}
	if len(got.Repositories) != len(want.Repositories) {
		t.Fatalf("len(Repositories): got %d, want %d", len(got.Repositories), len(want.Repositories))
	}
	for i, r := range want.Repositories {
		if got.Repositories[i] != r {
			t.Errorf("Repositories[%d]: got %+v, want %+v", i, got.Repositories[i], r)
		}
	}
}

func TestWriteConfigYAMLHeader(t *testing.T) {
	setupGritEnv(t)

	cfg, err := DefaultConfig()
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteConfig(cfg); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(data), "---\n") {
		t.Errorf("config file should start with '---\\n', got: %q", string(data))
	}
}

func TestWriteLoadConfigEmptyRepositories(t *testing.T) {
	setupGritEnv(t)

	if err := WriteConfig(Config{Root: "/tmp", IgnoreRoot: true}); err != nil {
		t.Fatal(err)
	}
	got, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	if len(got.Repositories) != 0 {
		t.Errorf("expected empty repositories after round-trip, got %d", len(got.Repositories))
	}
}

// --- AddRepoToConfig ---

func TestAddRepoToConfig_New(t *testing.T) {
	setupGritEnv(t)

	if err := AddRepoToConfig("myrepo", "myrepo"); err != nil {
		t.Fatal(err)
	}

	config, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(config.Repositories) != 1 {
		t.Fatalf("expected 1 repo, got %d", len(config.Repositories))
	}
	if config.Repositories[0].Name != "myrepo" {
		t.Errorf("Name = %q, want %q", config.Repositories[0].Name, "myrepo")
	}
	if config.Repositories[0].Path != "myrepo" {
		t.Errorf("Path = %q, want %q", config.Repositories[0].Path, "myrepo")
	}
}

func TestAddRepoToConfig_FirstEntry(t *testing.T) {
	setupGritEnv(t)

	initial, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(initial.Repositories) != 0 {
		t.Fatal("expected 0 repos initially")
	}

	if err := AddRepoToConfig("first", "first"); err != nil {
		t.Fatal(err)
	}

	after, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(after.Repositories) != 1 {
		t.Fatal("expected 1 repo after first add")
	}
}

func TestAddRepoToConfig_DuplicateName(t *testing.T) {
	setupGritEnv(t)

	if err := AddRepoToConfig("myrepo", "path-a"); err != nil {
		t.Fatal(err)
	}
	if err := AddRepoToConfig("myrepo", "path-b"); err != nil {
		t.Fatal(err)
	}

	config, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(config.Repositories) != 1 {
		t.Errorf("expected 1 repo after duplicate-name add, got %d", len(config.Repositories))
	}
}

func TestAddRepoToConfig_DuplicatePath(t *testing.T) {
	setupGritEnv(t)

	if err := AddRepoToConfig("name-a", "shared/path"); err != nil {
		t.Fatal(err)
	}
	if err := AddRepoToConfig("name-b", "shared/path"); err != nil {
		t.Fatal(err)
	}

	config, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(config.Repositories) != 1 {
		t.Errorf("expected 1 repo after duplicate-path add, got %d", len(config.Repositories))
	}
}

func TestAddRepoToConfig_MultipleDistinct(t *testing.T) {
	setupGritEnv(t)

	for _, name := range []string{"alpha", "beta", "gamma"} {
		if err := AddRepoToConfig(name, name); err != nil {
			t.Fatal(err)
		}
	}

	config, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if len(config.Repositories) != 3 {
		t.Errorf("expected 3 repos, got %d", len(config.Repositories))
	}
}

// --- RemoveRepoFromConfig ---

func mustAddRepo(t *testing.T, name, path string) {
	t.Helper()
	if err := AddRepoToConfig(name, path); err != nil {
		t.Fatal(err)
	}
}

func mustLoadConfig(t *testing.T) Config {
	t.Helper()
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func TestRemoveRepoFromConfig_Middle(t *testing.T) {
	setupGritEnv(t)

	mustAddRepo(t, "alpha", "alpha")
	mustAddRepo(t, "beta", "beta")
	mustAddRepo(t, "gamma", "gamma")

	if err := RemoveRepoFromConfig("beta"); err != nil {
		t.Fatal(err)
	}

	config := mustLoadConfig(t)
	if len(config.Repositories) != 2 {
		t.Fatalf("expected 2 repos after removal, got %d", len(config.Repositories))
	}
	for _, r := range config.Repositories {
		if r.Name == "beta" {
			t.Error("beta should have been removed")
		}
	}
}

func TestRemoveRepoFromConfig_First(t *testing.T) {
	setupGritEnv(t)

	mustAddRepo(t, "alpha", "alpha")
	mustAddRepo(t, "beta", "beta")

	if err := RemoveRepoFromConfig("alpha"); err != nil {
		t.Fatal(err)
	}

	config := mustLoadConfig(t)
	if len(config.Repositories) != 1 {
		t.Fatalf("expected 1 repo, got %d", len(config.Repositories))
	}
	if config.Repositories[0].Name != "beta" {
		t.Errorf("remaining repo should be beta, got %q", config.Repositories[0].Name)
	}
}

func TestRemoveRepoFromConfig_Last(t *testing.T) {
	setupGritEnv(t)

	mustAddRepo(t, "alpha", "alpha")
	mustAddRepo(t, "beta", "beta")

	if err := RemoveRepoFromConfig("beta"); err != nil {
		t.Fatal(err)
	}

	config := mustLoadConfig(t)
	if len(config.Repositories) != 1 {
		t.Fatalf("expected 1 repo, got %d", len(config.Repositories))
	}
	if config.Repositories[0].Name != "alpha" {
		t.Errorf("remaining repo should be alpha, got %q", config.Repositories[0].Name)
	}
}

func TestRemoveRepoFromConfig_OnlyEntry(t *testing.T) {
	setupGritEnv(t)

	mustAddRepo(t, "solo", "solo")
	if err := RemoveRepoFromConfig("solo"); err != nil {
		t.Fatal(err)
	}

	config := mustLoadConfig(t)
	if len(config.Repositories) != 0 {
		t.Errorf("expected empty repos after removing only entry, got %d", len(config.Repositories))
	}
}

func TestRemoveRepoFromConfig_NotFound(t *testing.T) {
	setupGritEnv(t)

	mustAddRepo(t, "alpha", "alpha")
	if err := RemoveRepoFromConfig("nonexistent"); err != nil {
		t.Fatal(err)
	}

	config := mustLoadConfig(t)
	if len(config.Repositories) != 1 {
		t.Errorf("config should be unchanged after removing nonexistent repo, got %d repos", len(config.Repositories))
	}
}

func TestRemoveRepoFromConfig_Wildcard(t *testing.T) {
	setupGritEnv(t)

	mustAddRepo(t, "gg-phoenix-api", "gg-phoenix-api")
	mustAddRepo(t, "gg-phoenix-web", "gg-phoenix-web")
	mustAddRepo(t, "other-repo", "other-repo")

	if err := RemoveRepoFromConfig("gg-phoenix-*"); err != nil {
		t.Fatal(err)
	}

	config := mustLoadConfig(t)
	if len(config.Repositories) != 1 {
		t.Fatalf("expected 1 repo after wildcard removal, got %d", len(config.Repositories))
	}
	if config.Repositories[0].Name != "other-repo" {
		t.Errorf("remaining repo should be other-repo, got %q", config.Repositories[0].Name)
	}
}

func TestRemoveRepoFromConfig_ShellSafeWildcard(t *testing.T) {
	setupGritEnv(t)

	mustAddRepo(t, "gg-phoenix-api", "gg-phoenix-api")
	mustAddRepo(t, "gg-phoenix-web", "gg-phoenix-web")
	mustAddRepo(t, "other-repo", "other-repo")

	if err := RemoveRepoFromConfig("gg-phoenix-%"); err != nil {
		t.Fatal(err)
	}

	config := mustLoadConfig(t)
	if len(config.Repositories) != 1 {
		t.Fatalf("expected 1 repo after %% wildcard removal, got %d", len(config.Repositories))
	}
	if config.Repositories[0].Name != "other-repo" {
		t.Errorf("remaining repo should be other-repo, got %q", config.Repositories[0].Name)
	}
}

func TestRemoveRepoFromConfig_InvalidPattern(t *testing.T) {
	setupGritEnv(t)

	mustAddRepo(t, "alpha", "alpha")
	if err := RemoveRepoFromConfig("["); err != nil {
		t.Fatal(err)
	}

	config := mustLoadConfig(t)
	if len(config.Repositories) != 1 {
		t.Errorf("config should be unchanged after invalid pattern, got %d repos", len(config.Repositories))
	}
}

// --- AppendHistory ---

func TestAppendHistory_Appends(t *testing.T) {
	setupGritEnv(t)

	if err := AppendHistory("grit pull"); err != nil {
		t.Fatal(err)
	}
	if err := AppendHistory("grit status"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(HistoryFile)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "grit pull") {
		t.Error("history should contain 'grit pull'")
	}
	if !strings.Contains(content, "grit status") {
		t.Error("history should contain 'grit status'")
	}
}

func TestAppendHistory_DoesNotOverwrite(t *testing.T) {
	setupGritEnv(t)

	if err := AppendHistory("first command"); err != nil {
		t.Fatal(err)
	}
	if err := AppendHistory("second command"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(HistoryFile)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "first command") || !strings.Contains(content, "second command") {
		t.Errorf("both entries should be present; got:\n%s", content)
	}
}

func TestAppendHistory_Format(t *testing.T) {
	setupGritEnv(t)

	if err := AppendHistory("grit pull"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(HistoryFile)
	if err != nil {
		t.Fatal(err)
	}
	line := strings.TrimSpace(string(data))
	if !strings.HasPrefix(line, "[") {
		t.Errorf("history entry should start with '[', got %q", line)
	}
	if !strings.Contains(line, "] grit pull") {
		t.Errorf("history entry should contain '] grit pull', got %q", line)
	}
}

func TestAppendHistory_CreatesFileIfNotExists(t *testing.T) {
	setupGritEnv(t)

	if err := os.Remove(HistoryFile); err != nil {
		t.Fatal(err)
	}

	if err := AppendHistory("test command"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(HistoryFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "test command") {
		t.Error("history file should have been created with the command entry")
	}
}

// --- AddAllRepos ---

func TestAddAllRepos_AddsGitDirs(t *testing.T) {
	setupGritEnv(t)

	makeFakeGitDir(t, "repo-a")
	makeFakeGitDir(t, "repo-b")
	if err := os.Mkdir("not-a-repo", 0755); err != nil {
		t.Fatal(err)
	}

	if err := AddAllRepos(); err != nil {
		t.Fatal(err)
	}

	names := repoNameSet(t)
	if !names["repo-a"] {
		t.Error("repo-a should have been added")
	}
	if !names["repo-b"] {
		t.Error("repo-b should have been added")
	}
	if names["not-a-repo"] {
		t.Error("not-a-repo should not have been added")
	}
}

func TestAddAllRepos_SkipsExistingRepos(t *testing.T) {
	setupGritEnv(t)

	makeFakeGitDir(t, "repo-a")
	mustAddRepo(t, "repo-a", "repo-a")

	if err := AddAllRepos(); err != nil {
		t.Fatal(err)
	}

	config := mustLoadConfig(t)
	count := 0
	for _, r := range config.Repositories {
		if r.Name == "repo-a" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("repo-a should appear exactly once, got %d", count)
	}
}

func TestAddAllRepos_SkipsGritDir(t *testing.T) {
	setupGritEnv(t)

	if err := AddAllRepos(); err != nil {
		t.Fatal(err)
	}

	config := mustLoadConfig(t)
	for _, r := range config.Repositories {
		if r.Name == GritDir || r.Path == GritDir {
			t.Errorf("grit dir %q should not be added as a repo", GritDir)
		}
	}
}

// --- RunGitCommandParallel ---

func mustWriteConfig(t *testing.T, cfg Config) {
	t.Helper()
	if err := WriteConfig(cfg); err != nil {
		t.Fatal(err)
	}
}

func TestRunGitCommandParallel_Completes(t *testing.T) {
	setupGritEnv(t)

	cwd := mustGetwd(t)
	makeRealGitRepo(t, "repo-a")
	makeRealGitRepo(t, "repo-b")

	mustWriteConfig(t, Config{
		Root: cwd,
		Repositories: []Repository{
			{Name: "repo-a", Path: "repo-a"},
			{Name: "repo-b", Path: "repo-b"},
		},
	})

	RunGitCommandParallel([]string{"status"})
}

func TestRunGitCommandParallel_EmptyRepos(t *testing.T) {
	setupGritEnv(t)

	// Empty config should return immediately without hanging.
	RunGitCommandParallel([]string{"status"})
}

func TestRunGitCommandParallel_SemaphoreLimitOne(t *testing.T) {
	setupGritEnv(t)

	cwd := mustGetwd(t)
	makeRealGitRepo(t, "repo-a")
	makeRealGitRepo(t, "repo-b")
	makeRealGitRepo(t, "repo-c")

	mustWriteConfig(t, Config{
		Root: cwd,
		Repositories: []Repository{
			{Name: "repo-a", Path: "repo-a"},
			{Name: "repo-b", Path: "repo-b"},
			{Name: "repo-c", Path: "repo-c"},
		},
	})

	t.Setenv("GRIT_MAX_CONCURRENT", "1")

	// Should complete without deadlock even with concurrency capped at 1.
	RunGitCommandParallel([]string{"status"})
}

func TestRunGitCommandParallel_InvalidSemaphoreValues(t *testing.T) {
	setupGritEnv(t)

	cwd := mustGetwd(t)
	makeRealGitRepo(t, "repo-a")
	mustWriteConfig(t, Config{
		Root:         cwd,
		Repositories: []Repository{{Name: "repo-a", Path: "repo-a"}},
	})

	// These values are all invalid or zero; the semaphore should be disabled
	// and the function should fall back to unbounded concurrency.
	for _, val := range []string{"abc", "-1", "0"} {
		t.Run("GRIT_MAX_CONCURRENT="+val, func(t *testing.T) {
			t.Setenv("GRIT_MAX_CONCURRENT", val)
			RunGitCommandParallel([]string{"status"})
		})
	}
}

func TestRunGitCommandParallel_ReportsFailedRepos(t *testing.T) {
	setupGritEnv(t)

	cwd := mustGetwd(t)
	makeRealGitRepo(t, "repo-ok")

	mustWriteConfig(t, Config{
		Root: cwd,
		Repositories: []Repository{
			{Name: "repo-ok", Path: "repo-ok"},
			{Name: "repo-missing", Path: "does-not-exist"},
		},
	})

	out := captureStdout(t, func() {
		RunGitCommandParallel([]string{"status"})
	})

	if !strings.Contains(out, "REPOSITORIES WITH ERRORS") {
		t.Fatalf("expected error summary in output, got:\n%s", out)
	}
	if !strings.Contains(out, "repo-missing:") {
		t.Fatalf("expected repo-missing in error summary, got:\n%s", out)
	}
	if !strings.Contains(out, "    ") {
		t.Fatalf("expected indented error lines in summary, got:\n%s", out)
	}
	if strings.Contains(out, "  repo-ok:") {
		t.Fatalf("repo-ok should not appear in error summary, got:\n%s", out)
	}
}

// --- RunCommandParallel ---

func TestRunCommandParallel_Completes(t *testing.T) {
	setupGritEnv(t)

	cwd := mustGetwd(t)
	if err := os.Mkdir("repo-a", 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir("repo-b", 0755); err != nil {
		t.Fatal(err)
	}

	mustWriteConfig(t, Config{
		Root: cwd,
		Repositories: []Repository{
			{Name: "repo-a", Path: "repo-a"},
			{Name: "repo-b", Path: "repo-b"},
		},
	})

	out := captureStdout(t, func() {
		RunCommandParallel("pwd", nil)
	})

	if !strings.Contains(out, "REPO-A") {
		t.Fatalf("expected REPO-A header in output, got:\n%s", out)
	}
	if !strings.Contains(out, "REPO-B") {
		t.Fatalf("expected REPO-B header in output, got:\n%s", out)
	}
	if !strings.Contains(out, "[pwd]") {
		t.Fatalf("expected [pwd] command display in output, got:\n%s", out)
	}
}

func TestRunCommandParallel_WithArgs(t *testing.T) {
	setupGritEnv(t)

	cwd := mustGetwd(t)
	if err := os.Mkdir("repo-a", 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join("repo-a", "hello.txt"), []byte("hi"), 0644); err != nil {
		t.Fatal(err)
	}

	mustWriteConfig(t, Config{
		Root:         cwd,
		Repositories: []Repository{{Name: "repo-a", Path: "repo-a"}},
	})

	out := captureStdout(t, func() {
		RunCommandParallel("ls", []string{"hello.txt"})
	})

	if !strings.Contains(out, "hello.txt") {
		t.Fatalf("expected hello.txt in output, got:\n%s", out)
	}
	if !strings.Contains(out, "[ls hello.txt]") {
		t.Fatalf("expected [ls hello.txt] command display, got:\n%s", out)
	}
}

func TestRunCommandParallel_ReportsFailedRepos(t *testing.T) {
	setupGritEnv(t)

	cwd := mustGetwd(t)
	if err := os.Mkdir("repo-ok", 0755); err != nil {
		t.Fatal(err)
	}

	mustWriteConfig(t, Config{
		Root: cwd,
		Repositories: []Repository{
			{Name: "repo-ok", Path: "repo-ok"},
			{Name: "repo-missing", Path: "does-not-exist"},
		},
	})

	out := captureStdout(t, func() {
		RunCommandParallel("pwd", nil)
	})

	if !strings.Contains(out, "REPOSITORIES WITH ERRORS") {
		t.Fatalf("expected error summary in output, got:\n%s", out)
	}
	if !strings.Contains(out, "repo-missing:") {
		t.Fatalf("expected repo-missing in error summary, got:\n%s", out)
	}
}

// --- helpers ---

func mustGetwd(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return cwd
}

func repoNameSet(t *testing.T) map[string]bool {
	t.Helper()
	names := make(map[string]bool)
	for _, r := range mustLoadConfig(t).Repositories {
		names[r.Name] = true
	}
	return names
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	fn()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = oldStdout

	var buf strings.Builder
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	if err := r.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}
