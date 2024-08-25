test "Fruit" {
  let first = Fruit::Apple
  let got = first.to_string()
  let want = "apple"
  assert_eq!(got, want)
  //
  let want =
    #|"apple"
  let got_parse = Fruit::parse!(want)
  assert_eq!(got_parse, first)
  //
  let mut threw_error = false
  let _ = try {
    Fruit::parse!("")
  } catch {
    _ => {
      threw_error = true
      Fruit::Apple
    }
  }
  assert_true!(threw_error)
  //
  let json_value = @jsonutil.to_json(first)
  let got = json_value.stringify()
  assert_eq!(got, want)
  //
  let got_parse = Fruit::from_json(json_value).unwrap()
  assert_eq!(got_parse, first)
  //
  let want_none = Fruit::from_json(Json::String(""))
  assert_eq!(want_none, None)
}

test "GhostGang" {
  let first = GhostGang::Blinky
  let got = first.to_string()
  let want = "blinky"
  assert_eq!(got, want)
  //
  let want =
    #|"blinky"
  let got_parse = GhostGang::parse!(want)
  assert_eq!(got_parse, first)
  //
  let mut threw_error = false
  let _ = try {
    GhostGang::parse!("")
  } catch {
    _ => {
      threw_error = true
      GhostGang::Blinky
    }
  }
  assert_true!(threw_error)
  //
  let json_value = @jsonutil.to_json(first)
  let got = json_value.stringify()
  assert_eq!(got, want)
  //
  let got_parse = GhostGang::from_json(json_value).unwrap()
  assert_eq!(got_parse, first)
  //
  let want_none = GhostGang::from_json(Json::String(""))
  assert_eq!(want_none, None)
}

test "ComplexObject" {
  let default_object = ComplexObject::new()
  let got = @jsonutil.to_json(default_object)
    |> @json.stringify()
  let want =
    #|{"ghost":"blinky","aBoolean":false,"aString":"","anInt":0}
  assert_eq!(got, want)
  //
  let got_parse = ComplexObject::parse!(want)
  assert_eq!(got_parse, default_object)
  //
  let required_fields : ComplexObject = {
    ghost: Blinky,
    a_boolean: true,
    a_string: "aString",
    an_int: 0,
    an_optional_date: None,
  }
  let got = @jsonutil.to_json(required_fields)
    |> @json.stringify()
  let want =
    #|{"ghost":"blinky","aBoolean":true,"aString":"aString","anInt":0}
  assert_eq!(got, want)
  //
  let got_parse = ComplexObject::parse!(want)
  assert_eq!(got_parse, required_fields)
  //
  let optional_fields : ComplexObject = {
    ..required_fields,
    an_optional_date: Some("anOptionalDate"),
  }
  let got = @jsonutil.to_json(optional_fields)
    |> @json.stringify()
  let want =
    #|{"ghost":"blinky","aBoolean":true,"aString":"aString","anInt":0,"anOptionalDate":"anOptionalDate"}
  assert_eq!(got, want)
  //
  let got_parse = ComplexObject::parse!(want)
  assert_eq!(got_parse, optional_fields)
}
