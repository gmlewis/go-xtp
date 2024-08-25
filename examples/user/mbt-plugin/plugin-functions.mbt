/// Exported: processUser
pub fn exported_process_user() -> Int {
  let user = match User::parse?(@host.input_string()) {
      Ok(user) => user
      Err(e) => {
        @host.set_error("unable to parse user: \{e}")
        return 1 // failure
      }
    }
    |> process_user()
  user.to_json() |> @host.output_json_value()
  return 0 // success
}
