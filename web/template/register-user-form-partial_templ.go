// Code generated by templ@v0.2.364 DO NOT EDIT.

package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	helpers "github.com/pmwals09/yobs/internal"
)

func RegisterUserForm(f helpers.FormData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_1 := templ.GetChildren(ctx)
		if var_1 == nil {
			var_1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<form action=\"/user/register\" method=\"POST\" class=\"mt-4\" hx-post=\"/user/register\" hx-swap=\"outerHTML\"><section class=\"grid grid-cols-2 gap-y-2\"><label for=\"username\" class=\"self-center\">")
		if err != nil {
			return err
		}
		var_2 := `User Name`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" id=\"username\" name=\"username\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(f.Values["username"]))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"> ")
		if err != nil {
			return err
		}
		var var_3 string = string(f.Errors["username"])
		_, err = templBuffer.WriteString(templ.EscapeString(var_3))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" <label for=\"email\" class=\"self-center\">")
		if err != nil {
			return err
		}
		var_4 := `Email Address`
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"text\" id=\"email\" name=\"email\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(f.Values["email"]))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"> ")
		if err != nil {
			return err
		}
		var var_5 string = string(f.Errors["email"])
		_, err = templBuffer.WriteString(templ.EscapeString(var_5))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" <label for=\"password\" class=\"self-center\">")
		if err != nil {
			return err
		}
		var_6 := `Password`
		_, err = templBuffer.WriteString(var_6)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"password\" id=\"password\" name=\"password\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(f.Values["password"]))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"> ")
		if err != nil {
			return err
		}
		var var_7 string = string(f.Errors["password"])
		_, err = templBuffer.WriteString(templ.EscapeString(var_7))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" <label for=\"password-repeat\" class=\"self-center\">")
		if err != nil {
			return err
		}
		var_8 := `Re-enter Password`
		_, err = templBuffer.WriteString(var_8)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><input type=\"password\" id=\"password-repeat\" name=\"password-repeat\" value=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(f.Values["passwordRepeat"]))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"> ")
		if err != nil {
			return err
		}
		var var_9 string = string(f.Errors["passwordRepeat"])
		_, err = templBuffer.WriteString(templ.EscapeString(var_9))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</section><section class=\"my-4\"><button type=\"submit\" class=\"bg-gray-400 px-4 py-2 rounded-full mx-auto block hover:cursor-pointer\">")
		if err != nil {
			return err
		}
		var_10 := `Submit`
		_, err = templBuffer.WriteString(var_10)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button></section> ")
		if err != nil {
			return err
		}
		var var_11 string = string(f.Errors["overall"])
		_, err = templBuffer.WriteString(templ.EscapeString(var_11))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</form>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
