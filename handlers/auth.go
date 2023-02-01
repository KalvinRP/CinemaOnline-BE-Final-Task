package handlers

import (
	"encoding/json"
	authdto "finaltask/dto/auth"
	dto "finaltask/dto/result"
	"finaltask/models"
	"finaltask/pkg/bcrypt"
	jwtToken "finaltask/pkg/jwt"
	repositories "finaltask/repositories"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"gopkg.in/gomail.v2"
)

type handlerAuth struct {
	AuthRepository repositories.AuthRepository
}

func HandlerAuth(AuthRepository repositories.AuthRepository) *handlerAuth {
	return &handlerAuth{AuthRepository}
}

func (h *handlerAuth) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(authdto.AuthRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	password, err := bcrypt.HashPass(request.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	role := "user"
	if strings.Contains(request.Email, "@cinemaonline.com") {
		role = "admin"
	}

	image := "https://res.cloudinary.com/dqnazzgq6/image/upload/v1674442268/CinemaOnline/pp_oixdhd.png"

	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: password,
		Role:     role,
		Image:    image,
	}

	usertosend := models.UserDetail{
		Name:  request.Name,
		Email: request.Email,
	}

	if h.AuthRepository.IsExist(user) {
		w.WriteHeader(http.StatusConflict)
		response := dto.ErrorResult{Code: http.StatusConflict, Message: "Email is Already Registered"}
		json.NewEncoder(w).Encode(response)
	} else {
		claims := jwt.MapClaims{}
		claims["name"] = user.Name
		claims["email"] = user.Email
		claims["password"] = user.Password
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix() // 2 minutes expired

		token, _ := jwtToken.GenerateToken(&claims)
		SendVerif(token, usertosend)

		w.WriteHeader(http.StatusOK)
		response := dto.ErrorResult{Code: http.StatusOK, Message: "Verification email sent!"}
		json.NewEncoder(w).Encode(response)
	}
}

func (h *handlerAuth) Verify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(authdto.VerifyRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Halo"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.Token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "You are not allowed to be here."}
		json.NewEncoder(w).Encode(response)
		return
	}

	claims, err := jwtToken.DecodeToken(request.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := dto.ErrorResult{Code: http.StatusUnauthorized, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	name := claims["name"].(string)
	email := claims["email"].(string)
	password := claims["password"].(string)
	passcrypt, err := bcrypt.HashPass(password)

	role := "user"
	if strings.Contains(email, "@cinemaonline.com") {
		role = "admin"
	}

	image := "https://res.cloudinary.com/dqnazzgq6/image/upload/v1674442268/CinemaOnline/pp_oixdhd.png"

	user := models.User{
		Name:     name,
		Email:    email,
		Password: passcrypt,
		Role:     role,
		Image:    image,
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	if h.AuthRepository.IsExist(user) {
		w.WriteHeader(http.StatusConflict)
		response := dto.ErrorResult{Code: http.StatusConflict, Message: "Email is Already Registered"}
		json.NewEncoder(w).Encode(response)
		return
	} else {
		h.AuthRepository.Register(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
		}
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponse(user)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerAuth) ForgetPass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(authdto.ForgetRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := h.AuthRepository.Login(request.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Email is not registered!"}
		json.NewEncoder(w).Encode(response)
		return
	} else {
		claims := jwt.MapClaims{}
		claims["name"] = user.Name
		claims["email"] = user.Email
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix() // 2 minutes expired

		token, _ := jwtToken.GenerateToken(&claims)
		SendReset(token, user.Email)

		w.WriteHeader(http.StatusOK)
		response := dto.ErrorResult{Code: http.StatusOK, Message: "Verification email sent!"}
		json.NewEncoder(w).Encode(response)
	}
}

func (h *handlerAuth) ResetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(authdto.ResetRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	claims, err := jwtToken.DecodeToken(request.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := dto.ErrorResult{Code: http.StatusUnauthorized, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	email := claims["email"].(string)
	passcrypt, err := bcrypt.HashPass(request.Password)

	user, err := h.AuthRepository.ResetPass(email, passcrypt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: user}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerAuth) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(authdto.LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	user := models.User{
		Email:    request.Email,
		Password: request.Password,
	}

	user, err := h.AuthRepository.Login(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	isValid := bcrypt.CheckHash(request.Password, user.Password)
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Wrong Email or Password!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	claims := jwt.MapClaims{}
	claims["role"] = user.Role
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix() // 2 hours expired

	token, errGenerateToken := jwtToken.GenerateToken(&claims)
	if errGenerateToken != nil {
		log.Println(errGenerateToken)
		fmt.Println("Unauthorize")
		return
	}

	loginResponse := authdto.AuthResponse{
		Name:  user.Name,
		Token: token,
		Role:  user.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	response := dto.SuccessResult{Code: http.StatusOK, Data: loginResponse}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerAuth) CheckAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	user, err := h.AuthRepository.CheckAuth(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix() // 2 hours expired

	token, errGenerateToken := jwtToken.GenerateToken(&claims)
	if errGenerateToken != nil {
		log.Println(errGenerateToken)
		fmt.Println("Unauthorize")
		return
	}

	authResponse := authdto.AuthResponse{
		Token: token,
		Role:  user.Role,
		Name:  user.Name,
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: authResponse}
	json.NewEncoder(w).Encode(response)
}

func SendVerif(token string, user models.UserDetail) {
	var CONFIG_SMTP_HOST = "smtp.gmail.com"
	var CONFIG_SMTP_PORT = 587
	var CONFIG_SENDER_NAME = os.Getenv("SENDER_NAME")
	var CONFIG_AUTH_EMAIL = os.Getenv("EMAIL_SYSTEM")
	var CONFIG_AUTH_PASSWORD = os.Getenv("PASSWORD_SYSTEM")
	var WEB = os.Getenv("URL")

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", user.Email)
	mailer.SetHeader("Subject", "Verify Registration")
	mailer.SetBody("text/html", fmt.Sprintf(`<!DOCTYPE html>
	<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
	<head>
		<meta charset="utf-8"> <!-- utf-8 works for most cases -->
		<meta name="viewport" content="width=device-width"> <!-- Forcing initial-scale shouldn't be necessary -->
		<meta http-equiv="X-UA-Compatible" content="IE=edge"> <!-- Use the latest (edge) version of IE rendering engine -->
		<meta name="x-apple-disable-message-reformatting">  <!-- Disable auto-scale in iOS 10 Mail entirely -->
		<title></title> <!-- The title tag shows in email notifications, like Android 4.4. -->
	
		<link href="https://fonts.googleapis.com/css?family=Lato:300,400,700" rel="stylesheet">
	
		<!-- CSS Reset : BEGIN -->
		<style>
	
			/* What it does: Remove spaces around the email design added by some email clients. */
			/* Beware: It can remove the padding / margin and add a background color to the compose a reply window. */
			html,
	body {
		margin: 0 auto !important;
		padding: 0 !important;
		height: 100%% !important;
		width: 100%% !important;
		background: #1D201F;
	}
	
	/* What it does: Stops email clients resizing small text. */
	* {
		-ms-text-size-adjust: 100%%;
		-webkit-text-size-adjust: 100%%;
	}
	
	/* What it does: Centers email on Android 4.4 */
	div[style*="margin: 16px 0"] {
		margin: 0 !important;
	}
	
	/* What it does: Stops Outlook from adding extra spacing to tables. */
	table,
	td {
		mso-table-lspace: 0pt !important;
		mso-table-rspace: 0pt !important;
	}
	
	/* What it does: Fixes webkit padding issue. */
	table {
		border-spacing: 0 !important;
		border-collapse: collapse !important;
		table-layout: fixed !important;
		margin: 0 auto !important;
	}
	
	/* What it does: Uses a better rendering method when resizing images in IE. */
	img {
		-ms-interpolation-mode:bicubic;
	}
	
	/* What it does: Prevents Windows 10 Mail from underlining links despite inline CSS. Styles for underlined links should be inline. */
	a {
		text-decoration: none;
	}
	
	/* What it does: A work-around for email clients meddling in triggered links. */
	*[x-apple-data-detectors],  /* iOS */
	.unstyle-auto-detected-links *,
	.aBn {
		border-bottom: 0 !important;
		cursor: default !important;
		color: inherit !important;
		text-decoration: none !important;
		font-size: inherit !important;
		font-family: inherit !important;
		font-weight: inherit !important;
		line-height: inherit !important;
	}
	
	/* What it does: Prevents Gmail from displaying a download button on large, non-linked images. */
	.a6S {
		display: none !important;
		opacity: 0.01 !important;
	}
	
	/* What it does: Prevents Gmail from changing the text color in conversation threads. */
	.im {
		color: inherit !important;
	}
	
	/* If the above doesn't work, add a .g-img class to any image in question. */
	img.g-img + div {
		display: none !important;
	}
	
	/* What it does: Removes right gutter in Gmail iOS app: https://github.com/TedGoas/Cerberus/issues/89  */
	/* Create one of these media queries for each additional viewport size you'd like to fix */
	
	/* iPhone 4, 4S, 5, 5S, 5C, and 5SE */
	@media only screen and (min-device-width: 320px) and (max-device-width: 374px) {
		u ~ div .email-container {
			min-width: 320px !important;
		}
	}
	/* iPhone 6, 6S, 7, 8, and X */
	@media only screen and (min-device-width: 375px) and (max-device-width: 413px) {
		u ~ div .email-container {
			min-width: 375px !important;
		}
	}
	/* iPhone 6+, 7+, and 8+ */
	@media only screen and (min-device-width: 414px) {
		u ~ div .email-container {
			min-width: 414px !important;
		}
	}
	
		</style>
	
		<!-- CSS Reset : END -->
	
		<!-- Progressive Enhancements : BEGIN -->
		<style>
	
			.primary{
		background: #30e3ca;
	}
	.bg_dark{
		background: #C58882;
	}
	.bg_light{
		background: #FF87AB;
	}
	.bg_black{
		background: #CD2E71;
	}
	.bg_dark{
		background: #1D201F;
	}
	.email-section{
		padding:2.5em;
	}
	
	/*BUTTON*/
	.btn{
		padding: 10px 15px;
		display: inline-block;
	}
	.btn.btn-primary{
		border-radius: 5px;
		background: #CD2E71;
		color: #000000;
		font-weight: bold;
	}
	.btn.btn-white{
		border-radius: 5px;
		background: #ffffff;
		color: #000000;
	}
	.btn.btn-white-outline{
		border-radius: 5px;
		background: transparent;
		border: 1px solid #fff;
		color: #fff;
	}
	.btn.btn-black-outline{
		border-radius: 0px;
		background: transparent;
		border: 2px solid #000;
		color: #000;
		font-weight: 700;
	}
	
	h1,h2,h3,h4,h5,h6{
		font-family: 'Lato', sans-serif;
		color: whitesmoke;
		margin-top: 0;
		font-weight: 400;
	}
	
	body{
		font-family: 'Lato', sans-serif;
		font-weight: 400;
		font-size: 15px;
		line-height: 1.8;
		color: rgba(0,0,0,.4);
	}
	
	a{
		color: #30e3ca;
	}
	
	table{
	}
	/*LOGO*/
	
	.logo h1{
		margin: 0;
	}
	.logo h1 a{
		color: #30e3ca;
		font-size: 24px;
		font-weight: 700;
		font-family: 'Lato', sans-serif;
	}
	
	/*HERO*/
	.hero{
		position: relative;
		z-index: 0;
	}
	
	.hero .text{
		color: whitesmoke;
	}
	.hero .text h2{
		color: whitesmoke;
		font-size: 40px;
		margin-bottom: 0;
		font-weight: 400;
		line-height: 1.4;
	}
	.hero .text h3{
		font-size: 24px;
		font-weight: 300;
	}
	.hero .text h2 span{
		font-weight: 600;
		color: #30e3ca;
	}
	
	
	/*HEADING SECTION*/
	.heading-section{
	}
	.heading-section h2{
		color: #000000;
		font-size: 28px;
		margin-top: 0;
		line-height: 1.4;
		font-weight: 400;
	}
	.heading-section .subheading{
		margin-bottom: 20px !important;
		display: inline-block;
		font-size: 13px;
		text-transform: uppercase;
		letter-spacing: 2px;
		color: rgba(0,0,0,.4);
		position: relative;
	}
	.heading-section .subheading::after{
		position: absolute;
		left: 0;
		right: 0;
		bottom: -10px;
		content: '';
		width: 100%%;
		height: 2px;
		background: #30e3ca;
		margin: 0 auto;
	}
	
	.heading-section-white{
		color: rgba(255,255,255,.8);
	}
	.heading-section-white h2{
		font-family: 
		line-height: 1;
		padding-bottom: 0;
	}
	.heading-section-white h2{
		color: #ffffff;
	}
	.heading-section-white .subheading{
		margin-bottom: 0;
		display: inline-block;
		font-size: 13px;
		text-transform: uppercase;
		letter-spacing: 2px;
		color: rgba(255,255,255,.4);
	}
	
	
	ul.social{
		padding: 0;
	}
	ul.social li{
		display: inline-block;
		margin-right: 10px;
	}
	
	/*FOOTER*/
	
	.footer{
		border-top: 1px solid rgba(0,0,0,.05);
		color: rgba(0,0,0,.5);
	}
	.footer .heading{
		color: #000;
		font-size: 20px;
	}
	.footer ul{
		margin: 0;
		padding: 0;
	}
	.footer ul li{
		list-style: none;
		margin-bottom: 10px;
	}
	.footer ul li a{
		color: rgba(0,0,0,1);
	}
	
	
	@media screen and (max-width: 500px) {
	
	
	}
	
	
		</style>
	
	
	</head>
	
	<body width="100%%" style="margin: 0; padding: 0 !important; mso-line-height-rule: exactly; background-color: #f1f1f1;">
		<center style="width: 100%%; background-color: #D1DEDE;">
		<div style="display: none; font-size: 1px;max-height: 0px; max-width: 0px; opacity: 0; overflow: hidden; mso-hide: all; font-family: sans-serif;">
		  &zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;
		</div>
		<div style="max-width: 600px; margin: 0 auto;" class="email-container">
			<!-- BEGIN BODY -->
		  <table align="center" role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="margin: auto;">
			  <tr>
			  <td valign="middle" class="hero bg_dark" style="padding: 3em 0 2em 0;">
				<img src="./Icon.svg" alt="" style="width: 300px; max-width: 600px; height: auto; margin: auto; display: block;">
			  </td>
			  </tr><!-- end tr -->
					<tr>
			  <td valign="middle" class="hero bg_dark" style="padding: 2em 0 4em 0;">
				<table>
					<tr>
						<td>
							<div class="text" style="padding: 0 2.5em; text-align: center;">
								<h2>Please verify your email</h2>
								<h3>Watch movies, series, and many more just by clicks!</h3>
								<h4 style="margin: 0; color: red; font-weight: bold;">IGNORE IF YOU HAVEN'T REGISTER TO OUR SERVICE</h4>
								<p><a href="%s/verify/%s" class="btn btn-primary">Click to Verify</a></p>
								<p>This email will be expired in 15 minutes after sent</p>
							</div>
						</td>
					</tr>
				</table>
			  </td>
			  </tr><!-- end tr -->
		  <!-- 1 Column Text + Button : END -->
		  </table>
		  <table align="center" role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="margin: auto;">
			  <tr>
			  <td valign="middle" class="bg_light footer email-section">
				<table>
					<tr>
					<td valign="top" width="33.333%%" style="padding-top: 20px;">
					  <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%">
						<tr>
						  <td style="text-align: left; padding-right: 10px;">
							  <h3 class="heading">About</h3>
							  <p>This email and other related objects are fictional for learning projects. Please do not attempt to do real transactions in corresponding website.</p>
						  </td>
						</tr>
					  </table>
					</td>
					<td valign="top" width="33.333%%" style="padding-top: 20px;">
					  <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%">
						<tr>
						  <td style="text-align: left; padding-left: 5px; padding-right: 5px;">
							  <h3 class="heading">Contact Info</h3>
							  <ul>
										<li><span class="text">Elang IV, Sawah Lama, Ciputat, Tangerang Selatan, Banten, Indonesia</span></li>
										<li><span class="text">https://dumbways.id/</span></a></li>
									  </ul>
						  </td>
						</tr>
					  </table>
					</td>
				  </tr>
				</table>
			  </td>
			</tr><!-- end: tr -->
			<tr>
			</tr>
		  </table>
	
		</div>
	  </center>
	</body>
	</html>`, WEB, token))

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Pesan terkirim!")
}

func SendReset(token, email string) {
	var CONFIG_SMTP_HOST = "smtp.gmail.com"
	var CONFIG_SMTP_PORT = 587
	var CONFIG_SENDER_NAME = os.Getenv("SENDER_NAME")
	var CONFIG_AUTH_EMAIL = os.Getenv("EMAIL_SYSTEM")
	var CONFIG_AUTH_PASSWORD = os.Getenv("PASSWORD_SYSTEM")
	var WEB = os.Getenv("URL")

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Verify Registration")
	mailer.SetBody("text/html", fmt.Sprintf(`<!DOCTYPE html>
	<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
	<head>
		<meta charset="utf-8"> <!-- utf-8 works for most cases -->
		<meta name="viewport" content="width=device-width"> <!-- Forcing initial-scale shouldn't be necessary -->
		<meta http-equiv="X-UA-Compatible" content="IE=edge"> <!-- Use the latest (edge) version of IE rendering engine -->
		<meta name="x-apple-disable-message-reformatting">  <!-- Disable auto-scale in iOS 10 Mail entirely -->
		<title></title> <!-- The title tag shows in email notifications, like Android 4.4. -->
	
		<link href="https://fonts.googleapis.com/css?family=Lato:300,400,700" rel="stylesheet">
	
		<!-- CSS Reset : BEGIN -->
		<style>
	
			/* What it does: Remove spaces around the email design added by some email clients. */
			/* Beware: It can remove the padding / margin and add a background color to the compose a reply window. */
			html,
	body {
		margin: 0 auto !important;
		padding: 0 !important;
		height: 100%% !important;
		width: 100%% !important;
		background: #1D201F;
	}
	
	/* What it does: Stops email clients resizing small text. */
	* {
		-ms-text-size-adjust: 100%%;
		-webkit-text-size-adjust: 100%%;
	}
	
	/* What it does: Centers email on Android 4.4 */
	div[style*="margin: 16px 0"] {
		margin: 0 !important;
	}
	
	/* What it does: Stops Outlook from adding extra spacing to tables. */
	table,
	td {
		mso-table-lspace: 0pt !important;
		mso-table-rspace: 0pt !important;
	}
	
	/* What it does: Fixes webkit padding issue. */
	table {
		border-spacing: 0 !important;
		border-collapse: collapse !important;
		table-layout: fixed !important;
		margin: 0 auto !important;
	}
	
	/* What it does: Uses a better rendering method when resizing images in IE. */
	img {
		-ms-interpolation-mode:bicubic;
	}
	
	/* What it does: Prevents Windows 10 Mail from underlining links despite inline CSS. Styles for underlined links should be inline. */
	a {
		text-decoration: none;
	}
	
	/* What it does: A work-around for email clients meddling in triggered links. */
	*[x-apple-data-detectors],  /* iOS */
	.unstyle-auto-detected-links *,
	.aBn {
		border-bottom: 0 !important;
		cursor: default !important;
		color: inherit !important;
		text-decoration: none !important;
		font-size: inherit !important;
		font-family: inherit !important;
		font-weight: inherit !important;
		line-height: inherit !important;
	}
	
	/* What it does: Prevents Gmail from displaying a download button on large, non-linked images. */
	.a6S {
		display: none !important;
		opacity: 0.01 !important;
	}
	
	/* What it does: Prevents Gmail from changing the text color in conversation threads. */
	.im {
		color: inherit !important;
	}
	
	/* If the above doesn't work, add a .g-img class to any image in question. */
	img.g-img + div {
		display: none !important;
	}
	
	/* What it does: Removes right gutter in Gmail iOS app: https://github.com/TedGoas/Cerberus/issues/89  */
	/* Create one of these media queries for each additional viewport size you'd like to fix */
	
	/* iPhone 4, 4S, 5, 5S, 5C, and 5SE */
	@media only screen and (min-device-width: 320px) and (max-device-width: 374px) {
		u ~ div .email-container {
			min-width: 320px !important;
		}
	}
	/* iPhone 6, 6S, 7, 8, and X */
	@media only screen and (min-device-width: 375px) and (max-device-width: 413px) {
		u ~ div .email-container {
			min-width: 375px !important;
		}
	}
	/* iPhone 6+, 7+, and 8+ */
	@media only screen and (min-device-width: 414px) {
		u ~ div .email-container {
			min-width: 414px !important;
		}
	}
	
		</style>
	
		<!-- CSS Reset : END -->
	
		<!-- Progressive Enhancements : BEGIN -->
		<style>
	
			.primary{
		background: #30e3ca;
	}
	.bg_dark{
		background: #C58882;
	}
	.bg_light{
		background: #FF87AB;
	}
	.bg_black{
		background: #CD2E71;
	}
	.bg_dark{
		background: #1D201F;
	}
	.email-section{
		padding:2.5em;
	}
	
	/*BUTTON*/
	.btn{
		padding: 10px 15px;
		display: inline-block;
	}
	.btn.btn-primary{
		border-radius: 5px;
		background: #CD2E71;
		color: #000000;
		font-weight: bold;
	}
	.btn.btn-white{
		border-radius: 5px;
		background: #ffffff;
		color: #000000;
	}
	.btn.btn-white-outline{
		border-radius: 5px;
		background: transparent;
		border: 1px solid #fff;
		color: #fff;
	}
	.btn.btn-black-outline{
		border-radius: 0px;
		background: transparent;
		border: 2px solid #000;
		color: #000;
		font-weight: 700;
	}
	
	h1,h2,h3,h4,h5,h6{
		font-family: 'Lato', sans-serif;
		color: whitesmoke;
		margin-top: 0;
		font-weight: 400;
	}
	
	body{
		font-family: 'Lato', sans-serif;
		font-weight: 400;
		font-size: 15px;
		line-height: 1.8;
		color: rgba(0,0,0,.4);
	}
	
	a{
		color: #30e3ca;
	}
	
	table{
	}
	/*LOGO*/
	
	.logo h1{
		margin: 0;
	}
	.logo h1 a{
		color: #30e3ca;
		font-size: 24px;
		font-weight: 700;
		font-family: 'Lato', sans-serif;
	}
	
	/*HERO*/
	.hero{
		position: relative;
		z-index: 0;
	}
	
	.hero .text{
		color: whitesmoke;
	}
	.hero .text h2{
		color: whitesmoke;
		font-size: 40px;
		margin-bottom: 0;
		font-weight: 400;
		line-height: 1.4;
	}
	.hero .text h3{
		font-size: 24px;
		font-weight: 300;
	}
	.hero .text h2 span{
		font-weight: 600;
		color: #30e3ca;
	}
	
	
	/*HEADING SECTION*/
	.heading-section{
	}
	.heading-section h2{
		color: #000000;
		font-size: 28px;
		margin-top: 0;
		line-height: 1.4;
		font-weight: 400;
	}
	.heading-section .subheading{
		margin-bottom: 20px !important;
		display: inline-block;
		font-size: 13px;
		text-transform: uppercase;
		letter-spacing: 2px;
		color: rgba(0,0,0,.4);
		position: relative;
	}
	.heading-section .subheading::after{
		position: absolute;
		left: 0;
		right: 0;
		bottom: -10px;
		content: '';
		width: 100%%;
		height: 2px;
		background: #30e3ca;
		margin: 0 auto;
	}
	
	.heading-section-white{
		color: rgba(255,255,255,.8);
	}
	.heading-section-white h2{
		font-family: 
		line-height: 1;
		padding-bottom: 0;
	}
	.heading-section-white h2{
		color: #ffffff;
	}
	.heading-section-white .subheading{
		margin-bottom: 0;
		display: inline-block;
		font-size: 13px;
		text-transform: uppercase;
		letter-spacing: 2px;
		color: rgba(255,255,255,.4);
	}
	
	
	ul.social{
		padding: 0;
	}
	ul.social li{
		display: inline-block;
		margin-right: 10px;
	}
	
	/*FOOTER*/
	
	.footer{
		border-top: 1px solid rgba(0,0,0,.05);
		color: rgba(0,0,0,.5);
	}
	.footer .heading{
		color: #000;
		font-size: 20px;
	}
	.footer ul{
		margin: 0;
		padding: 0;
	}
	.footer ul li{
		list-style: none;
		margin-bottom: 10px;
	}
	.footer ul li a{
		color: rgba(0,0,0,1);
	}
	
	
	@media screen and (max-width: 500px) {
	
	
	}
	
	
		</style>
	
	
	</head>
	
	<body width="100%%" style="margin: 0; padding: 0 !important; mso-line-height-rule: exactly; background-color: #f1f1f1;">
		<center style="width: 100%%; background-color: #D1DEDE;">
		<div style="display: none; font-size: 1px;max-height: 0px; max-width: 0px; opacity: 0; overflow: hidden; mso-hide: all; font-family: sans-serif;">
		  &zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;&zwnj;&nbsp;
		</div>
		<div style="max-width: 600px; margin: 0 auto;" class="email-container">
			<!-- BEGIN BODY -->
		  <table align="center" role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="margin: auto;">
			  <tr>
			  <td valign="middle" class="hero bg_dark" style="padding: 3em 0 2em 0;">
				<img src="./Icon.svg" alt="" style="width: 300px; max-width: 600px; height: auto; margin: auto; display: block;">
			  </td>
			  </tr><!-- end tr -->
					<tr>
			  <td valign="middle" class="hero bg_dark" style="padding: 2em 0 4em 0;">
				<table>
					<tr>
						<td>
							<div class="text" style="padding: 0 2.5em; text-align: center;">
								<h2>Reset Password Request</h2>
								<h3>Click the button below to reset your password.</h3>
								<h4 style="margin: 0; color: red; font-weight: bold;">IGNORE IF YOU THINK THIS IS A MISTAKE</h4>
								<p><a href="%s/reset/%s" class="btn btn-primary">Click to Reset</a></p>
								<p>This email will be expired in 15 minutes after sent</p>
							</div>
						</td>
					</tr>
				</table>
			  </td>
			  </tr><!-- end tr -->
		  <!-- 1 Column Text + Button : END -->
		  </table>
		  <table align="center" role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="margin: auto;">
			  <tr>
			  <td valign="middle" class="bg_light footer email-section">
				<table>
					<tr>
					<td valign="top" width="33.333%%" style="padding-top: 20px;">
					  <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%">
						<tr>
						  <td style="text-align: left; padding-right: 10px;">
							  <h3 class="heading">About</h3>
							  <p>This email and other related objects are fictional for learning projects. Please do not attempt to do real transactions in corresponding website.</p>
						  </td>
						</tr>
					  </table>
					</td>
					<td valign="top" width="33.333%%" style="padding-top: 20px;">
					  <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%">
						<tr>
						  <td style="text-align: left; padding-left: 5px; padding-right: 5px;">
							  <h3 class="heading">Contact Info</h3>
							  <ul>
										<li><span class="text">Elang IV, Sawah Lama, Ciputat, Tangerang Selatan, Banten, Indonesia</span></li>
										<li><span class="text">https://dumbways.id/</span></a></li>
									  </ul>
						  </td>
						</tr>
					  </table>
					</td>
				  </tr>
				</table>
			  </td>
			</tr><!-- end: tr -->
			<tr>
			</tr>
		  </table>
	
		</div>
	  </center>
	</body>
	</html>`, WEB, token))

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Pesan terkirim!")
}
