-- hada-commerce harness init
-- Loads agents, skills, registers team template, and wires custom tools.

claudio.agents.load_dir(PROJECT_CLAUDIO_DIR .. "/agents")
claudio.skills.load_dir(PROJECT_CLAUDIO_DIR .. "/skills")

-- ─── Team Template ────────────────────────────────────────────────────────────

claudio.teams.register_template({
  name        = "hada-team",
  description = "Full hada-commerce team: Go backend, harness tools, React frontend, reviewer, and devops.",
  members = {
    { subagent_type = "investigator",   model = "claude-haiku-4-5-20251001" },
    { subagent_type = "go-backend",     model = "claude-sonnet-4-6"         },
    { subagent_type = "go-agent-tools", model = "claude-sonnet-4-6"         },
    { subagent_type = "react-frontend", model = "claude-sonnet-4-6"         },
    { subagent_type = "reviewer",       model = "claude-sonnet-4-6"         },
    { subagent_type = "devops",         model = "claude-haiku-4-5-20251001" },
  },
})

-- ─── Custom Tools ─────────────────────────────────────────────────────────────

-- Build: compile the Go backend and report errors
claudio.tools.register({
  name        = "Build",
  description = [[
    Compile the hada-commerce Go backend. Runs 'go build ./...' and 'go vet ./...' from
    the backend/ directory. Returns compiler output and any vet warnings.
    Use this after backend changes to verify the code compiles before spawning reviewer.
  ]],
  agents = { "principal", "reviewer", "go-backend", "go-agent-tools" },
  schema = [[{"type": "object", "properties": {}}]],
  execute = function(_input)
    local handle = io.popen("cd " .. PROJECT_CLAUDIO_DIR .. "/../backend && go build ./... 2>&1 && go vet ./... 2>&1")
    local result = handle:read("*a")
    local ok, _, code = handle:close()
    if code ~= 0 then
      return "BUILD FAILED (exit " .. tostring(code) .. "):\n" .. result
    end
    if result == "" then
      return "BUILD OK — no errors."
    end
    return "BUILD OK (vet warnings):\n" .. result
  end,
})

-- Test: run all Go tests
claudio.tools.register({
  name        = "Test",
  description = [[
    Run all Go backend tests with 'go test ./...' from the backend/ directory.
    Returns test output including PASS/FAIL per package. Use after Build to verify correctness.
  ]],
  agents = { "principal", "reviewer" },
  schema = [[{"type": "object", "properties": {}}]],
  execute = function(_input)
    local handle = io.popen("cd " .. PROJECT_CLAUDIO_DIR .. "/../backend && go test ./... 2>&1")
    local result = handle:read("*a")
    local ok, _, code = handle:close()
    if code ~= 0 then
      return "TESTS FAILED (exit " .. tostring(code) .. "):\n" .. result
    end
    return "TESTS PASSED:\n" .. result
  end,
})

-- ServiceStatus: check docker-compose service health
claudio.tools.register({
  name        = "ServiceStatus",
  description = [[
    Show the status of docker-compose services (postgres, redis).
    Returns 'docker compose ps' output. Use to verify infra is running before debugging DB issues.
  ]],
  agents = { "principal", "devops", "reviewer" },
  schema = [[{"type": "object", "properties": {}}]],
  execute = function(_input)
    local handle = io.popen("cd " .. PROJECT_CLAUDIO_DIR .. "/.. && docker compose ps 2>&1")
    local result = handle:read("*a")
    handle:close()
    return result
  end,
})

-- Migrate: apply pending SQL migrations in order
claudio.tools.register({
  name        = "Migrate",
  description = [[
    Apply SQL migrations from backend/migrations/ in sequential order using psql.
    Requires DATABASE_URL environment variable to be set.
    Only runs .up.sql files. Migrations are append-only — never re-runs already-applied ones.
    IMPORTANT: Confirm with the user before running against a production database.
  ]],
  agents = { "principal", "devops" },
  schema = [[{
    "type": "object",
    "required": ["database_url"],
    "properties": {
      "database_url": {
        "type": "string",
        "description": "PostgreSQL connection string, e.g. postgres://hada:hada@localhost:5433/hada"
      }
    }
  }]],
  execute = function(input)
    local db_url = input.database_url or ""
    -- Validate: must look like a postgres URL
    if not db_url:match("^postgres") then
      return "Error: database_url must start with 'postgres' — got: " .. db_url
    end
    local migrations_dir = PROJECT_CLAUDIO_DIR .. "/../backend/migrations"
    local handle = io.popen("ls " .. migrations_dir .. "/*.up.sql 2>&1 | sort")
    local files = handle:read("*a")
    handle:close()

    if files == "" then
      return "No migration files found in " .. migrations_dir
    end

    local results = "Running migrations:\n"
    for file in files:gmatch("[^\n]+") do
      local mh = io.popen("psql '" .. db_url .. "' -f '" .. file .. "' 2>&1")
      local out = mh:read("*a")
      local _, _, code = mh:close()
      results = results .. "\n[" .. (code == 0 and "OK" or "FAIL") .. "] " .. file .. "\n" .. out
    end
    return results
  end,
})

