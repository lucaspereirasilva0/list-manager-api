# AI Agent Guidelines for Testing

## Overview

This document provides specific guidelines for AI agents when writing or modifying tests in this project. All tests MUST follow the patterns defined in CLAUDE.md.

## Table-Driven Test Checklist

Before writing any test, ask yourself:

- [ ] Are there multiple test cases with **similar structure**?
- [ ] Do the cases differ only in **inputs and expected outputs**?
- [ ] Are the assertions **consistent** across all cases?

**If YES** → Use table-driven test pattern

**If NO** → Consider separate test functions

## Mandatory Table-Driven Pattern

```go
func Test<FunctionName>(t *testing.T) {
	tests := []struct {
		name              string
		given<FieldName>  <Type>
		want<FieldName>    <Type>
		wantErr           error
	}{
		{
			name:              "Given_<Condition>_When_<Action>_Then_<Result>",
			given<FieldName>:  <Value>,
			want<FieldName>:    <ExpectedValue>,
			wantErr:           nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mock := setupMock()
			mock.On("Method", tt.given<FieldName>).Return(tt.want<FieldName>, tt.wantErr)

			// Act
			result, err := FunctionUnderTest(tt.given<FieldName>)

			// Assert
			require.Equal(t, tt.want<FieldName>, result)
			if tt.wantErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
```

## Common Anti-Patterns (AVOID)

### ❌ Anti-Pattern 1: Multiple `t.Run` with duplicate code

```go
// ✗ AVOID
func TestFunction(t *testing.T) {
	t.Run("Given_Case1_When_Then", func(t *testing.T) {
		mock := new(Mock)
		mock.On("Method").Return("result1")
		result := Function()
		require.Equal(t, "result1", result)
	})

	t.Run("Given_Case2_When_Then", func(t *testing.T) {
		mock := new(Mock)  // DUPLICATE CODE
		mock.On("Method").Return("result2")  // DUPLICATE CODE
		result := Function()  // DUPLICATE CODE
		require.Equal(t, "result2", result)  // DUPLICATE CODE
	})
}
```

### ❌ Anti-Pattern 2: Conditional logic inside table-driven test

```go
// ✗ AVOID
for _, tt := range tests {
	t.Run(tt.name, func(t *testing.T) {
		// ... setup and act ...

		// AVOID THIS
		if tt.name == "Given_SpecialCase" {
			require.Contains(t, result, "special-field")
			return
		}

		require.Equal(t, tt.value, result)
	})
}
```

### ❌ Anti-Pattern 3: Non-BDD naming

```go
// ✗ AVOID
t.Run("TestHappyPath", func(t *testing.T) { /* ... */ })
t.Run("TestErrorCase", func(t *testing.T) { /* ... */ })
t.Run("ShouldFail", func(t *testing.T) { /* ... */ })

// ✓ USE
t.Run("Given_ValidInput_When_Processing_Then_Success", func(t *testing.T) { /* ... */ })
t.Run("Given_InvalidInput_When_Processing_Then_ReturnsError", func(t *testing.T) { /* ... */ })
```

## When to Separate Tests

Create separate test functions when:

1. **Different aspects are being validated**
   ```go
   func TestFunction(t *testing.T) {
       // Table-driven for status values
   }

   func TestFunctionFields(t *testing.T) {
       // Separate test for field validation
   }
   ```

2. **Test structure is fundamentally different**
   - One validates return values
   - Other validates side effects or metrics

3. **Complex conditional logic would be needed** in table-driven loop

## Test Structure Template

Every test function MUST follow this structure:

```go
func Test<FeatureName>(t *testing.T) {
	// 1. Define test cases
	tests := []struct {
		name              string  // BDD format: Given_When_Then
		given<Input>      <Type>
		want<Output>      <Type>
		wantErr           error
	}{/* ... */}

	// 2. Iterate through cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 3. Arrange: Setup mocks, dependencies
			// 4. Act: Execute function
			// 5. Assert: Verify results
		})
	}
}
```

## Field Naming Convention

Test struct fields MUST use these prefixes:

- `given` - Input values or initial state
- `want` - Expected output values
- `got` - For debugging (rarely used in production)

```go
tests := []struct {
	name              string
	givenUserID       string    // ✓ CORRECT
	givenName         string    // ✓ CORRECT
	wantHTTPStatus    int       // ✓ CORRECT
	wantResponse      Response  // ✓ CORRECT
	wantErr           error     // ✓ CORRECT
}{/* ... */}

// ✗ AVOID
tests := []struct {
	name         string
	inputID      string    // ✗ WRONG: use 'givenID'
	expectedResp Response  // ✗ WRONG: use 'wantResponse'
	error        error     // ✗ WRONG: use 'wantErr'
}{/* ... */}
```

## AI Agent Checklist

Before submitting any test code, verify:

- [ ] Uses table-driven pattern when appropriate
- [ ] Follows BDD naming: `Given_When_Then`
- [ ] No duplicate code between test cases
- [ ] Uses `given*` and `want*` field prefixes
- [ ] Uses `require` for critical assertions
- [ ] No conditional logic inside table-driven loop
- [ ] Separate tests for different validation logic
- [ ] Mock helper functions following `mock*` pattern
- [ ] All tests are independent (no shared state)

## Examples in This Project

Reference existing test files for patterns:
- `cmd/api/handlers/handlers_test.go` - Table-driven examples
- `cmd/api/handlers/health_test.go` - Correct pattern
- `internal/service/service_test.go` - Table-driven with mocks

## Quality Gates

All tests must pass before submission:
```bash
go test ./... -v
go test ./<package> -cover
```

Coverage goal: >80% for all packages
