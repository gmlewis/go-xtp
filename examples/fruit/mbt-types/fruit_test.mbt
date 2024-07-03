test "ComplexObject" {
  let default_object = ComplexObject::new()
  let got = @jsonutil.to_json(default_object)
    |> @jsonutil.stringify(spaces=0, newline=false)
  let want =
    #|{"ghost":"blinky","aBoolean":false,"aString":"","anInt":0}
  @test.eq(got, want)!
  //
  let required_fields : ComplexObject = {
    ghost: Blinky,
    a_boolean: true,
    a_string: "aString",
    an_int: 0,
    an_optional_date: None,
  }
  let got = @jsonutil.to_json(required_fields)
    |> @jsonutil.stringify(spaces=0, newline=false)
  let want =
    #|{"ghost":"blinky","aBoolean":true,"aString":"aString","anInt":0}
  @test.eq(got, want)!
  //
  let optional_fields : ComplexObject = {
    ..required_fields,
    an_optional_date: Some("anOptionalDate"),
  }
  let got = @jsonutil.to_json(optional_fields)
    |> @jsonutil.stringify(spaces=0, newline=false)
  let want =
    #|{"ghost":"blinky","aBoolean":true,"aString":"aString","anInt":0,"anOptionalDate":"anOptionalDate"}
  @test.eq(got, want)!
}
