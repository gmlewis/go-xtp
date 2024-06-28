test "Address" {
  let default_object = Address::new()
  let got = default_object.to_json()
  let want =
    #|{"street":"street"}`
  @assertion.assert_eq(got, want)?
  //
  let required_fields : Address = { street: "street" }
  let got = required_fields.to_json()
  let want =
    #|{"street":"street"}`
  @assertion.assert_eq(got, want)?
}

test "User" {
  let default_object = User::new()
  let got = default_object.to_json()
  let want =
    #|{"street":""}`
  @assertion.assert_eq(got, want)?
  //
  let required_fields : User = { age: None, email: None, address: None }
  let got = required_fields.to_json()
  let want =
    #|{}`
  @assertion.assert_eq(got, want)?
  //
  let optional_fields : User = {
    ..required_fields,
    age: Some(42),
    email: Some("email@example.com"),
    address: Some(Address::new()),
  }
  let got = optional_fields.to_json()
  let want =
    #|{"age":42,"email":"email@example.com","address":{"street":""}}`
  @assertion.assert_eq(got, want)?
}
