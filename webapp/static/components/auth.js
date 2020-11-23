const login = {
  data() {
    return {
      username: '',
      password: '',
      rememberme: false
    }
  },
  template: `
<div class='content' @keyup.enter='login'>
  <header>
    <h3 class='d-flex justify-content-center align-items-center' style='height: 100%'>Log In</h3>
  </header>
  <div class='login'>
    <div class='form-group'>
      <label for='username'>Username</label>
      <input autofocus class='form-control' v-model.trim='username' id='username' maxlength=20 placeholder='Username' required>
    </div>
    <div class='form-group'>
      <label for='password'>Password</label>
      <input class='form-control' type='password' v-model.trim='password' id='password' maxlength=20 placeholder='Password' required>
    </div>
    <div class='form-group form-check'>
      <input type='checkbox' class='form-check-input' v-model='rememberme' id='rememberme'>
      <label class='form-check-label' for='rememberme'>Remember Me</label>
    </div>
    <hr>
    <button class='btn btn-primary login' @click='login'>Log In</button>
  </div>
</div>`,
  mounted() { document.title = 'Log In' },
  methods: {
    login() {
      if (!username.checkValidity())
        BootstrapButtons.fire('Error', 'Username cannot be empty.', 'error')
      else if (!password.checkValidity())
        BootstrapButtons.fire('Error', 'Password cannot be empty.', 'error')
      else post('/login', {
        username: this.username,
        password: this.password,
        rememberme: this.rememberme
      }).then(resp => {
        if (!resp.ok) resp.text().then(err =>
          BootstrapButtons.fire('Error', err, 'error'))
        else window.location = '/'
      })
    }
  }
}

const setting = {
  data() {
    return {
      password: '',
      password1: '',
      password2: '',
      validated: false
    }
  },
  template: `
<div class='content' @keyup.enter='setting'>
  <header style='padding-left: 20px'>
    <h3>Setting</h3>
    <hr>
  </header>
  <div style='margin-left:120px; width:250px;' :class="{ 'was-validated': validated }">
    <div class='form-group'>
      <label for='password'>Current Password</label>
      <input class='form-control' type='password' v-model.trim='password' id='password' maxlength=20 required>
      <div class='invalid-feedback'>This field is required.</div>
    </div>
    <div class='form-group'>
      <label for='password1'>New Password</label>
      <input class='form-control' type='password' v-model.trim='password1' id='password1' maxlength=20 required>
      <div class='invalid-feedback'>This field is required.</div>
    </div>
    <div class='form-group'>
      <label for='password2'>Confirm Password</label>
      <input class='form-control' type='password' v-model.trim='password2' id='password2' maxlength=20 required>
      <div class='invalid-feedback'>This field is required.</div>
      <small class='form-text text-muted'>Max password length: 20 characters.</small>
    </div>
    <button class='btn btn-primary' @click='setting'>Change</button>
    <button class='btn btn-primary' @click='goback()'>Cancel</button>
  </div>
</div>`,
  mounted() {
    document.title = 'Setting'
    window.addEventListener('keyup', this.cancel)
  },
  beforeUnmount: function () { window.removeEventListener('keyup', this.cancel) },
  methods: {
    setting() {
      if (valid()) {
        this.validated = false
        post('/setting', {
          password: this.password,
          password1: this.password1,
          password2: this.password2
        }).then(resp => {
          if (!resp.ok) resp.text().then(err =>
            BootstrapButtons.fire('Error', err, 'error'))
          else resp.json().then(json => {
            if (json.status == 1)
              BootstrapButtons.fire('Success', 'Your password has changed. Please Re-login!', 'success')
                .then(() => window.location = '/')
            else
              BootstrapButtons.fire('Error', json.message, 'error')
                .then(() => {
                  if (json.error == 1) this.password = ''
                  else {
                    this.password1 = ''
                    this.password2 = ''
                  }
                })
          })
        })
      }
      else this.validated = true
    },
    goback() { this.$router.go(-1) },
    cancel(event) { if (event.key == 'Escape') this.goback() }
  }
}
