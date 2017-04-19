package interpreter

import (
	"testing"

	"github.com/goruby/goruby/object"
)

func TestInterpreterInterpret(t *testing.T) {
	t.Run("return proper result", func(t *testing.T) {
		input := `
			def foo
				3
			end

			x = 5

			def add x, y
				x + y
			end

			add foo, x
			`
		i := New()

		out, err := i.Interpret(input)
		if err != nil {
			panic(err)
		}

		res, ok := out.(*object.Integer)
		if !ok {
			t.Logf("Expected *object.Integer, got %T\n", out)
			t.Fail()
		}

		if res.Value != 8 {
			t.Logf("Expected result to equal 8, got %d\n", res.Value)
			t.Fail()
		}
	})
	t.Run("return proper result with changed env", func(t *testing.T) {
		input := `
			def foo
				3
			end

			def add x, y
				x + y
			end

			add foo, x
			`
		env := object.NewMainEnvironment()
		env.Set("x", &object.Integer{Value: 3})
		i := New()
		i.SetEnvironment(env)

		out, err := i.Interpret(input)
		if err != nil {
			panic(err)
		}

		res, ok := out.(*object.Integer)
		if !ok {
			t.Logf("Expected *object.Integer, got %T\n", out)
			t.Fail()
		}

		if res.Value != 6 {
			t.Logf("Expected result to equal 8, got %d\n", res.Value)
			t.Fail()
		}
	})
}

func TestInterpretClasses(t *testing.T) {
	t.Run("evaluate class object", func(t *testing.T) {
		input := `
		class Foo
			def foo
				3
			end
		end
		Foo
		`

		i := New()

		evaluated, err := i.Interpret(input)
		if err != nil {
			t.Logf("Expected no error, got %T:%v", err, err)
			t.FailNow()
		}

		classObject, ok := evaluated.(object.RubyClassObject)
		if !ok {
			t.Logf("Expected class object, got %T", evaluated)
			t.FailNow()
		}

		_, ok = classObject.Methods().Get("foo")
		if !ok {
			t.Logf("Expected class object to have method foo")
			t.Fail()
		}
	})
	t.Run("self after class definition", func(t *testing.T) {
		input := `
		class Foo
			def foo
				3
			end
		end
		self
		`

		i := New()

		evaluated, err := i.Interpret(input)
		if err != nil {
			t.Logf("Expected no error, got %T:%v", err, err)
			t.FailNow()
		}

		self, ok := evaluated.(*object.Self)
		if !ok {
			t.Logf("Expected self, got %T", evaluated)
			t.FailNow()
		}

		if self.Name != "main" {
			t.Logf("Expected main object, got %s", self.Name)
			t.Fail()
		}
	})
}
