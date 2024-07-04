/// Exported: processUser
pub fn exported_process_user() -> Int {
  User::parse(@host.input_string())!!.unwrap()
  |> process_user()
  |> @jsonutil.to_json()
  |> @host.output_json_value()
  return 0 // success
}
