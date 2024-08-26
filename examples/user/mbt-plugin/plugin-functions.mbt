/// Exported: processUser
pub fn exported_process_user() -> Int {
  let input = @host.input_string()
  let user : User = match @json.from_json?(input.to_json()) {
    Ok(user) => user
    _ => {
      @host.set_error("unable to parse input \{input}")
      return 1 // failure
    }
  }
  process_user(user).to_json() |> @host.output_json_value()
  return 0 // success
}
