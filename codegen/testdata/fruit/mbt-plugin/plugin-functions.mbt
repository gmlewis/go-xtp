/// Exported: voidFunc
pub fn exported_void_func() -> Int {
  void_func()
  return 0 // success
}

/// Exported: primitiveTypeFunc
pub fn exported_primitive_type_func() -> Int {
  let result = @json.parse?(@host.input_string())
  let input = match result {
    Ok(String(s)) => s
    _ => {
      @host.set_error(
        "exported_primitive_type_func: unable to parse JSON input",
      )
      return 1 // failure
    }
  }
  let output = primitive_type_func(input).to_json()
  @host.output_json_value(output)
  return 0 // success
}

/// Exported: referenceTypeFunc
pub fn exported_reference_type_func() -> Int {
  let input = @host.input_string()
  let fruit : Fruit = match @json.from_json?(input.to_json()) {
    Ok(fruit) => fruit
    _ => {
      @host.set_error("unable to parse input \{input}")
      return 1 // failure
    }
  }
  reference_type_func(fruit).to_json() |> @host.output_json_value()
  return 0 // success
}
