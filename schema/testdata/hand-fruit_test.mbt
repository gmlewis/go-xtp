fn assert_eq[T : Show + Eq](lhs : T, rhs : T) -> Unit {
  if lhs != rhs {
    abort("assert_eq failed.\n    lhs: \(lhs)\n    rhs: \(rhs)")
  }
}

test {
  let requiredFields = {
    ghost: Blinky,
    aBoolean: true,
    aString: "aString",
    anInt: 0,
    anOptionalDate: None,
  }
  let got = requiredFields.to_json()
  let want =
    #|{"ghost":"blinky","aBoolean":true,"aString":"aString","anInt":0}`
  assert_eq(got, want)
  //
  let optionalFields = {
    ..requiredFields,
    anOptionalDate: Some("anOptionalDate"),
  }
  let got = optionalFields.to_json()
  let want =
    #|{"ghost":"blinky","aBoolean":true,"aString":"aString","anInt":0}`
  assert_eq(got, want)
}
