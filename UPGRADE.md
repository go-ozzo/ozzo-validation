# Upgrade Instructions

## Upgrade from 3.x to 4.x
* We deprecated the `NewStrnigRule`  function in favor of the `NewStringValidator`. The new function 
gets `validator function`, `code`(to use as the translation key in error translation) and `message` as
its params. The following snippet shows how to modify your code if you want to define new string rule:
 ```go
// 3.x
// Assume Email is your custom rule:
var Email = validation.NewStringRule(govalidator.IsEmail, "must be a valid email address")

// 4.x
var Email = validation.NewStringValidator(govalidator.IsEmail, "email","must be a valid email address")
```

## Upgrade from 2.x to 3.x

* Instead of using `StructRules` to define struct validation rules, use `ValidateStruct()` to declare and perform
  struct validation. The following code snippet shows how to modify your code:
```go
// 2.x usage
err := validation.StructRules{}.
	Add("Street", validation.Required, validation.Length(5, 50)).
	Add("City", validation.Required, validation.Length(5, 50)).
	Add("State", validation.Required, validation.Match(regexp.MustCompile("^[A-Z]{2}$"))).
	Add("Zip", validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))).
	Validate(a)

// 3.x usage
err := validation.ValidateStruct(&a,
	validation.Field(&a.Street, validation.Required, validation.Length(5, 50)),
	validation.Field(&a.City, validation.Required, validation.Length(5, 50)),
	validation.Field(&a.State, validation.Required, validation.Match(regexp.MustCompile("^[A-Z]{2}$"))),
	validation.Field(&a.Zip, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
)
```

* Instead of using `Rules` to declare a rule list and use it to validate a value, call `Validate()` with the rules directly.
```go
data := "example"

// 2.x usage
rules := validation.Rules{
	validation.Required,      
	validation.Length(5, 100),
	is.URL,                   
}
err := rules.Validate(data)

// 3.x usage
err := validation.Validate(data,
	validation.Required,      
	validation.Length(5, 100),
	is.URL,                   
)
```

* The default struct tags used for determining error keys is changed from `validation` to `json`. You may modify
  `validation.ErrorTag` to change it back.