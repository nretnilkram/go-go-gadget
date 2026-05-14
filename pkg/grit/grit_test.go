package grit

import (
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
	WriteConfig(DefaultConfig())

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

	config := DefaultConfig()

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
	WriteConfig(want)
	got := LoadConfig()

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

	WriteConfig(DefaultConfig())

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

	WriteConfig(Config{Root: "/tmp", IgnoreRoot: true})
	got := LoadConfig()

	if len(got.Repositories) != 0 {
		t.Errorf("expected empty repositories after round-trip, got %d", len(got.Repositories))
	}
}

// --- AddRepoToConfig ---

func TestAddRepoToConfig_New(t *testing.T) {
	setupGritEnv(t)

	AddRepoToConfig("myrepo", "myrepo")

	config := LoadConfig()
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

	if len(LoadConfig().Repositories) != 0 {
		t.Fatal("expected 0 repos initially")
	}

	AddRepoToConfig("first", "first")

	if len(LoadConfig().Repositories) != 1 {
		t.Fatal("expected 1 repo after first add")
	}
}

func TestAddRepoToConfig_DuplicateName(t *testing.T) {
	setupGritEnv(t)

	AddRepoToConfig("myrepo", "path-a")
	AddRepoToConfig("myrepo", "path-b") // same name, different path

	config := LoadConfig()
	if len(config.Repositories) != 1 {
		t.Errorf("expected 1 repo after duplicate-name add, got %d", len(config.Repositories))
	}
}

func TestAddRepoToConfig_DuplicatePath(t *testing.T) {
	setupGritEnv(t)

	AddRepoToConfig("name-a", "shared/path")
	AddRepoToConfig("name-b", "shared/path") // same path, different name

	config := LoadConfig()
	if len(config.Repositories) != 1 {
		t.Errorf("expected 1 repo after duplicate-path add, got %d", len(config.Repositories))
	}
}

func TestAddRepoToConfig_MultipleDistinct(t *testing.T) {
	setupGritEnv(t)

	AddRepoToConfig("alpha", "alpha")
	AddRepoToConfig("beta", "beta")
	AddRepoToConfig("gamma", "gamma")

	config := LoadConfig()
	if len(config.Repositories) != 3 {
		t.Errorf("expected 3 repos, got %d", len(config.Repositories))
	}
}

// --- RemoveRepoFromConfig ---

func TestRemoveRepoFromConfig_Middle(t *testing.T) {
	setupGritEnv(t)

	AddRepoToConfig("alpha", "alpha")
	AddRepoToConfig("beta", "beta")
	AddRepoToConfig("gamma", "gamma")

	RemoveRepoFromConfig("beta")

	config := LoadConfig()
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

	AddRepoToConfig("alpha", "alpha")
	AddRepoToConfig("beta", "beta")

	RemoveRepoFromConfig("alpha")

	config := LoadConfig()
	if len(config.Repositories) != 1 {
		t.Fatalf("expected 1 repo, got %d", len(config.Repositories))
	}
	if config.Repositories[0].Name != "beta" {
		t.Errorf("remaining repo should be beta, got %q", config.Repositories[0].Name)
	}
}

func TestRemoveRepoFromConfig_Last(t *testing.T) {
	setupGritEnv(t)

	AddRepoToConfig("alpha", "alpha")
	AddRepoToConfig("beta", "beta")

	RemoveRepoFromConfig("beta")

	config := LoadConfig()
	if len(config.Repositories) != 1 {
		t.Fatalf("expected 1 repo, got %d", len(config.Repositories))
	}
	if config.Repositories[0].Name != "alpha" {
		t.Errorf("remaining repo should be alpha, got %q", config.Repositories[0].Name)
	}
}

func TestRemoveRepoFromConfig_OnlyEntry(t *testing.T) {
	setupGritEnv(t)

	AddRepoToConfig("solo", "solo")
	RemoveRepoFromConfig("solo")

	config := LoadConfig()
	if len(config.Repositories) != 0 {
		t.Errorf("expected empty repos after removing only entry, got %d", len(config.Repositories))
	}
}

func TestRemoveRepoFromConfig_NotFound(t *testing.T) {
	setupGritEnv(t)

	AddRepoToConfig("alpha", "alpha")
	RemoveRepoFromConfig("nonexistent")

	config := LoadConfig()
	if len(config.Repositories) != 1 {
		t.Errorf("config should be unchanged after removing nonexistent repo, got %d repos", len(config.Repositories))
	}
}

// --- AppendHistory ---

func TestAppendHistory_Appends(t *testing.T) {
	setupGritEnv(t)

	AppendHistory("grit pull")
	AppendHistory("grit status")

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

	AppendHistory("first command")
	AppendHistory("second command")

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

	AppendHistory("grit pull")

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

	// O_CREATE flag should re-create the file.
	AppendHistory("test command")

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

	AddAllRepos()

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
	AddRepoToConfig("repo-a", "repo-a")

	AddAllRepos()

	count := 0
	for _, r := range LoadConfig().Repositories {
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

	AddAllRepos()

	for _, r := range LoadConfig().Repositories {
		if r.Name == GritDir || r.Path == GritDir {
			t.Errorf("grit dir %q should not be added as a repo", GritDir)
		}
	}
}

// --- RunGitCommandParallel ---

func TestRunGitCommandParallel_Completes(t *testing.T) {
	setupGritEnv(t)

	cwd := mustGetwd(t)
	makeRealGitRepo(t, "repo-a")
	makeRealGitRepo(t, "repo-b")

	WriteConfig(Config{
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

	WriteConfig(Config{
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
	WriteConfig(Config{
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

// --- RunGitCommandSynchronous ---

func TestRunGitCommandSynchronous_Completes(t *testing.T) {
	setupGritEnv(t)

	cwd := mustGetwd(t)
	makeRealGitRepo(t, "repo-a")
	makeRealGitRepo(t, "repo-b")

	WriteConfig(Config{
		Root: cwd,
		Repositories: []Repository{
			{Name: "repo-a", Path: "repo-a"},
			{Name: "repo-b", Path: "repo-b"},
		},
	})

	RunGitCommandSynchronous([]string{"status"})
}

func TestRunGitCommandSynchronous_EmptyRepos(t *testing.T) {
	setupGritEnv(t)

	RunGitCommandSynchronous([]string{"status"})
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
	for _, r := range LoadConfig().Repositories {
		names[r.Name] = true
	}
	return names
}
