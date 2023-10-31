// Code generated by templ@v0.2.364 DO NOT EDIT.

package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/user"
)

func SignupPage(user *user.User, f helpers.FormData) templ.Component {
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
		var_2 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			templBuffer, templIsBuffer := w.(*bytes.Buffer)
			if !templIsBuffer {
				templBuffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templBuffer)
			}
			_, err = templBuffer.WriteString("<h1 class=\"text-center text-[2em]\">")
			if err != nil {
				return err
			}
			var_3 := `Join us!`
			_, err = templBuffer.WriteString(var_3)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h1> <section class=\"mb-4 w-1/2\"><form action=\"/user/register\" method=\"POST\" class=\"mt-4\" hx-post=\"/user/register\" hx-swap=\"outerHTML\" hx-select=\"#register-user-form\" id=\"register-user-form\"><section class=\"grid grid-cols-2 gap-y-2\"><label for=\"username\" class=\"self-center\">")
			if err != nil {
				return err
			}
			var_4 := `User Name`
			_, err = templBuffer.WriteString(var_4)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><div><input type=\"text\" id=\"username\" name=\"username\" value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(f.Values["username"]))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			err = RenderStandardError(f.Errors["username"]).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div><label for=\"email\" class=\"self-center\">")
			if err != nil {
				return err
			}
			var_5 := `Email Address`
			_, err = templBuffer.WriteString(var_5)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><div><input type=\"text\" id=\"email\" name=\"email\" value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(f.Values["email"]))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			err = RenderStandardError(f.Errors["email"]).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div><label for=\"password\" class=\"self-center\">")
			if err != nil {
				return err
			}
			var_6 := `Password`
			_, err = templBuffer.WriteString(var_6)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><div><input type=\"password\" id=\"password\" name=\"password\" value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(f.Values["password"]))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			err = RenderStandardError(f.Errors["password"]).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div><label for=\"password-repeat\" class=\"self-center\">")
			if err != nil {
				return err
			}
			var_7 := `Re-enter Password`
			_, err = templBuffer.WriteString(var_7)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><div><input type=\"password\" id=\"password-repeat\" name=\"password-repeat\" value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(f.Values["passwordRepeat"]))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			err = RenderStandardError(f.Errors["passwordRepeat"]).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div></section><section class=\"my-4\"><button type=\"submit\" class=\"bg-gray-400 px-4 py-2 rounded-full mx-auto block hover:cursor-pointer\">")
			if err != nil {
				return err
			}
			var_8 := `Submit`
			_, err = templBuffer.WriteString(var_8)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</button></section>")
			if err != nil {
				return err
			}
			err = RenderStandardError(f.Errors["overall"]).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</form></section>")
			if err != nil {
				return err
			}
			if !templIsBuffer {
				_, err = io.Copy(w, templBuffer)
			}
			return err
		})
		err = base(user).Render(templ.WithChildren(ctx, var_2), templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
