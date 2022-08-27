package stdlib

func RunArray() CommandFunc {
  set := map[string]CommandFunc{}
  return makeEnsemble("array", set)
}
