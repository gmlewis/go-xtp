test "ComplexObject" {
  let default_object = ComplexObject::new()
  let got = default_object.to_json()
  let want =
    #|{"ghost":"blinky","aBoolean":false,"aString":"","anInt":0}
  @assertion.assert_eq(got, want)?
  //
  let required_fields : ComplexObject = {
    ghost: Blinky,
    a_boolean: true,
    a_string: "aString",
    an_int: 0,
    an_optional_date: None,
  }
  let got = required_fields.to_json()
  let want =
    #|{"ghost":"blinky","aBoolean":true,"aString":"aString","anInt":0}
  @assertion.assert_eq(got, want)?
  //
  let optional_fields : ComplexObject = {
    ..required_fields,
    an_optional_date: Some("anOptionalDate"),
  }
  let got = optional_fields.to_json()
  let want =
    #|{"ghost":"blinky","aBoolean":true,"aString":"aString","anInt":0,"anOptionalDate":"anOptionalDate"}
  @assertion.assert_eq(got, want)?
}
