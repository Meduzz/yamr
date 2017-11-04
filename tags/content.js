riot.tag2('content', '<register if="{controller == ⁗register⁗}"></register> <profile if="{controller == ⁗profile⁗}"></profile> <search if="{controller == ⁗home⁗}"></search>', '', '', function(opts) {
    this.controller = ""

    let r = riot.route.create()
    const self = this

    r("/register", register)
    r("/profile", profile)
    r("/", home)

    function register() {
      self.controller = "register"
      self.update()
    }

    function profile() {
      self.controller = "profile"
      self.update()
    }

    function home() {
      self.controller = "home"
      self.update()
    }
});
