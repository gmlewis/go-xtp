test "Fruit" {
  let first = Fruit::Apple
  let got = first.to_string()
  let want = "apple"
  @test.eq(got, want)!
  //
  let got_parse = Fruit::parse(want)!
  @test.eq(got_parse, first)!
  //
  let mut threw_error = false
  let _ = try {
    Fruit::parse("")!
  } catch {
    _ => {
      threw_error = true
      Fruit::Apple
    }
  }
  @test.is_true(threw_error)!
  //
  let got_parse = Fruit::from_json(@json.JsonValue::String(want)).unwrap()
  @test.eq(got_parse, first)!
  //
  let want_none = Fruit::from_json(@json.JsonValue::String(""))
  @test.eq(want_none, None)!
}

test "GhostGang" {
  let first = GhostGang::Blinky
  let got = first.to_string()
  let want = "blinky"
  @test.eq(got, want)!
  //
  let got_parse = GhostGang::parse(want)!
  @test.eq(got_parse, first)!
  //
  let mut threw_error = false
  let _ = try {
    GhostGang::parse("")!
  } catch {
    _ => {
      threw_error = true
      GhostGang::Blinky
    }
  }
  @test.is_true(threw_error)!
  //
  let got_parse = GhostGang::from_json(@json.JsonValue::String(want)).unwrap()
  @test.eq(got_parse, first)!
  //
  let want_none = GhostGang::from_json(@json.JsonValue::String(""))
  @test.eq(want_none, None)!
}

test "ComplexObject" {
  let default_object = ComplexObject::new()
  let got = @jsonutil.to_json(default_object)
    |> @jsonutil.stringify(spaces=0, newline=false)
  let want =
    #|{"ghost":"blinky","aBoolean":false,"aString":"","anInt":0}
  @test.eq(got, want)!
  //
  let got_parse = ComplexObject::parse(want)!
  @test.eq(got_parse, default_object)!
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
  let got_parse = ComplexObject::parse(want)!
  @test.eq(got_parse, required_fields)!
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
  //
  let got_parse = ComplexObject::parse(want)!
  @test.eq(got_parse, optional_fields)!
}
