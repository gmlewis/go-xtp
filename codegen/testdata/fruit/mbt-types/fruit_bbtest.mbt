test "Fruit.to_string() works as expected" {
  let first = Fruit::Apple
  let got = first.to_string()
  let want = "apple"
  assert_eq!(got, want)
}

test "Fruit::from_json() works as expected" {
  let got_parse : Fruit = @json.from_json!("apple".to_json())
  let want = Fruit::Apple
  assert_eq!(got_parse, want)
  //
  let mut threw_error = false
  let _ = try {
    @json.from_json!("")
  } catch {
    _ => {
      threw_error = true
      Fruit::Apple
    }
  }
  assert_true!(threw_error)
}

test "Fruit::to_json() works as expected" {
  let first = Fruit::Apple
  let json_value = first.to_json()
  let got = json_value.stringify(escape_slash=false)
  let want =
    #|"apple"
  assert_eq!(got, want)
  //
  let got_parse : Fruit = @json.from_json!(@json.parse!(want))
  assert_eq!(got_parse, first)
}

test "GhostGang.to_string() works as expected" {
  let first = GhostGang::Blinky
  let got = first.to_string()
  let want = "blinky"
  assert_eq!(got, want)
}

test "GhostGang::from_json() works as expected" {
  let got_parse : GhostGang = @json.from_json!("blinky".to_json())
  let want = GhostGang::Blinky
  assert_eq!(got_parse, want)
  //
  let mut threw_error = false
  let _ = try {
    @json.from_json!("")
  } catch {
    _ => {
      threw_error = true
      GhostGang::Blinky
    }
  }
  assert_true!(threw_error)
}

test "GhostGang::to_json() works as expected" {
  let first = GhostGang::Blinky
  let json_value = first.to_json()
  let got = json_value.stringify(escape_slash=false)
  let want =
    #|"blinky"
  assert_eq!(got, want)
  //
  let got_parse : GhostGang = @json.from_json!(@json.parse!(want))
  assert_eq!(got_parse, first)
}

test "ComplexObject.to_json and .from_json work as expected on default object" {
  let default_object = ComplexObject::new()
  let got = default_object.to_json().stringify(escape_slash=false)
  let want =
    #|{"ghost":"blinky","aBoolean":false,"aString":"","anInt":0}
  assert_eq!(got, want)
  //
  let got_parse : ComplexObject = @json.from_json!(@json.parse!(want))
  assert_eq!(got_parse, default_object)
}

test "ComplexObject.to_json and .from_json work as expected on object only containing required fields" {
  let required_fields : ComplexObject = {
    ghost: Blinky,
    a_boolean: true,
    a_string: "aString",
    an_int: 0,
    an_optional_date: None,
  }
  let got = required_fields.to_json() |> @json.stringify(escape_slash=false)
  let want =
    #|{"ghost":"blinky","aBoolean":true,"aString":"aString","anInt":0}
  assert_eq!(got, want)
  //
  let got_parse : ComplexObject = @json.from_json!(@json.parse!(want))
  assert_eq!(got_parse, required_fields)
}

test "ComplexObject.to_json and .from_json work as expected on object with optional fields" {
  let optional_fields : ComplexObject = {
    ..ComplexObject::new(),
    an_optional_date: Some("anOptionalDate"),
  }
  let got = optional_fields.to_json() |> @json.stringify(escape_slash=false)
  let want =
    #|{"ghost":"blinky","aBoolean":false,"aString":"","anInt":0,"anOptionalDate":"anOptionalDate"}
  assert_eq!(got, want)
  //
  let got_parse : ComplexObject = @json.from_json!(@json.parse!(want))
  assert_eq!(got_parse, optional_fields)
}
