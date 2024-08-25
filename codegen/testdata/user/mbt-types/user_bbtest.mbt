test "Address" {
  let default_object = Address::new()
  let got = @jsonutil.to_json(default_object)
    |> @json.stringify()
  let want =
    #|{"street":""}
  assert_eq!(got, want)!
  //
  let got_parse = Address::parse(want)!
  assert_eq!(got_parse, default_object)!
  //
  let required_fields : Address = {
    street: "street",
  }
  let got = @jsonutil.to_json(required_fields)
    |> @json.stringify()
  let want =
    #|{"street":"street"}
  assert_eq!(got, want)!
  //
  let got_parse = Address::parse(want)!
  assert_eq!(got_parse, required_fields)!
}

test "User" {
  let default_object = User::new()
  let got = @jsonutil.to_json(default_object)
    |> @json.stringify()
  let want =
    #|{}
  assert_eq!(got, want)!
  //
  let got_parse = User::parse(want)!
  assert_eq!(got_parse, default_object)!
  //
  let required_fields : User = {
    age: None,
    email: None,
    address: None,
  }
  let got = @jsonutil.to_json(required_fields)
    |> @json.stringify()
  let want =
    #|{}
  assert_eq!(got, want)!
  //
  let got_parse = User::parse(want)!
  assert_eq!(got_parse, required_fields)!
  //
  let optional_fields : User = {
    ..required_fields,
    age: Some(0),
    email: Some("email"),
    address: Some(Address::new()),
  }
  let got = @jsonutil.to_json(optional_fields)
    |> @json.stringify()
  let want =
    #|{"age":0,"email":"email","address":{"street":""}}
  assert_eq!(got, want)!
  //
  let got_parse = User::parse(want)!
  assert_eq!(got_parse, optional_fields)!
}
