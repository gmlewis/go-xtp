test "Address" {
  let default_object = Address::new()
  let got = @jsonutil.to_json(default_object)
    |> @jsonutil.stringify(spaces=0, newline=false)
  let want =
    #|{"street":""}
  @assertion.assert_eq(got, want)?
  //
  let got_parse = Address::parse(want)!
  @assertion.assert_eq(got_parse, default_object)?
  //
  let required_fields : Address = {
    street: "street",
  }
  let got = @jsonutil.to_json(required_fields)
    |> @jsonutil.stringify(spaces=0, newline=false)
  let want =
    #|{"street":"street"}
  @assertion.assert_eq(got, want)?
  //
  let got_parse = Address::parse(want)!
  @assertion.assert_eq(got_parse, required_fields)?
}

test "User" {
  let default_object = User::new()
  let got = @jsonutil.to_json(default_object)
    |> @jsonutil.stringify(spaces=0, newline=false)
  let want =
    #|{}
  @assertion.assert_eq(got, want)?
  //
  let got_parse = User::parse(want)!
  @assertion.assert_eq(got_parse, default_object)?
  //
  let required_fields : User = {
    age: None,
    email: None,
    address: None,
  }
  let got = @jsonutil.to_json(required_fields)
    |> @jsonutil.stringify(spaces=0, newline=false)
  let want =
    #|{}
  @assertion.assert_eq(got, want)?
  //
  let got_parse = User::parse(want)!
  @assertion.assert_eq(got_parse, required_fields)?
  //
  let optional_fields : User = {
    ..required_fields,
    age: Some(0),
    email: Some("email"),
    address: Some(Address::new()),
  }
  let got = @jsonutil.to_json(optional_fields)
    |> @jsonutil.stringify(spaces=0, newline=false)
  let want =
    #|{"age":0,"email":"email","address":{"street":""}}
  @assertion.assert_eq(got, want)?
  //
  let got_parse = User::parse(want)!
  @assertion.assert_eq(got_parse, optional_fields)?
}
