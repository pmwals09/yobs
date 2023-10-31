// Code generated by templ@v0.2.364 DO NOT EDIT.

package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/opportunity"
	"github.com/pmwals09/yobs/internal/models/user"
)

func HomePage(user *user.User, opportunities []opportunity.Opportunity, f helpers.FormData) templ.Component {
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
			var_3 := `Always look ahead`
			_, err = templBuffer.WriteString(var_3)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h1> <section class=\"mb-4\" id=\"opportunity-list\"><h2 class=\"text-[1.5em]\">")
			if err != nil {
				return err
			}
			var_4 := `Opportunities`
			_, err = templBuffer.WriteString(var_4)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h2><table class=\"table-auto border-collapse w-full\" id=\"main-content\"><thead><tr><th class=\"border-b text-left\">")
			if err != nil {
				return err
			}
			var_5 := `Company Name`
			_, err = templBuffer.WriteString(var_5)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</th><th class=\"border-b text-left\">")
			if err != nil {
				return err
			}
			var_6 := `Role`
			_, err = templBuffer.WriteString(var_6)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</th><th class=\"border-b text-left\">")
			if err != nil {
				return err
			}
			var_7 := `URL`
			_, err = templBuffer.WriteString(var_7)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</th><th class=\"border-b text-left\">")
			if err != nil {
				return err
			}
			var_8 := `Status`
			_, err = templBuffer.WriteString(var_8)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</th><th class=\"border-b text-left\">")
			if err != nil {
				return err
			}
			var_9 := `Application Date`
			_, err = templBuffer.WriteString(var_9)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</th><th class=\"border-b text-left\">")
			if err != nil {
				return err
			}
			var_10 := `Actions`
			_, err = templBuffer.WriteString(var_10)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</th></tr></thead><tbody>")
			if err != nil {
				return err
			}
			for _, o := range opportunities {
				_, err = templBuffer.WriteString("<tr class=\"hover:bg-sky-100\"><td class=\"border-b\">")
				if err != nil {
					return err
				}
				var var_11 string = o.CompanyName
				_, err = templBuffer.WriteString(templ.EscapeString(var_11))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</td><td class=\"border-b\">")
				if err != nil {
					return err
				}
				var var_12 string = o.Role
				_, err = templBuffer.WriteString(templ.EscapeString(var_12))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</td><td class=\"border-b\">")
				if err != nil {
					return err
				}
				var var_13 string = o.URL
				_, err = templBuffer.WriteString(templ.EscapeString(var_13))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</td><td class=\"border-b\">")
				if err != nil {
					return err
				}
				var var_14 string = string(o.Status)
				_, err = templBuffer.WriteString(templ.EscapeString(var_14))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</td><td class=\"border-b\">")
				if err != nil {
					return err
				}
				var var_15 string = formatApplicationDate(o.ApplicationDate)
				_, err = templBuffer.WriteString(templ.EscapeString(var_15))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</td><td class=\"border-b\"><a class=\"btn\" href=\"")
				if err != nil {
					return err
				}
				var var_16 templ.SafeURL = insertIDIntoHref("opportunities/{}", o.ID)
				_, err = templBuffer.WriteString(templ.EscapeString(string(var_16)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\">")
				if err != nil {
					return err
				}
				var_17 := `Details`
				_, err = templBuffer.WriteString(var_17)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</a><!--")
				if err != nil {
					return err
				}
				var_18 := ` <a class="btn" href={ insertIDIntoHref("opportunities/{}/edit", o.ID) }>Edit</a> `
				_, err = templBuffer.WriteString(var_18)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("--></td></tr>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</tbody></table></section> <section class=\"mb-4\"><h2 class=\"text-[1.5em]\">")
			if err != nil {
				return err
			}
			var_19 := `New Opportunity`
			_, err = templBuffer.WriteString(var_19)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h2><form method=\"POST\" action=\"/opportunities\" hx-post=\"/opportunities\" hx-target=\"#opportunity-list\" hx-select=\"#opportunity-list\" id=\"opportunity-form\" class=\"w-1/2\"><section class=\"grid grid-cols-2 gap-y-2\"><label for=\"opportunity-name\">")
			if err != nil {
				return err
			}
			var_20 := `Company Name`
			_, err = templBuffer.WriteString(var_20)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"opportunity-name\" value=\"\"><label for=\"opportunity-role\">")
			if err != nil {
				return err
			}
			var_21 := `Role`
			_, err = templBuffer.WriteString(var_21)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input type=\"text\" name=\"opportunity-role\" value=\"\"><label for=\"opportunity-url\">")
			if err != nil {
				return err
			}
			var_22 := `URL`
			_, err = templBuffer.WriteString(var_22)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input type=\"url\" name=\"opportunity-url\" value=\"\"><label for=\"opportunity-description\">")
			if err != nil {
				return err
			}
			var_23 := `Description`
			_, err = templBuffer.WriteString(var_23)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><textarea rows=\"3\" name=\"opportunity-description\"></textarea><label for=\"opportunity-date\">")
			if err != nil {
				return err
			}
			var_24 := `Application Date (if applicable)`
			_, err = templBuffer.WriteString(var_24)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><input type=\"date\" name=\"opportunity-date\" value=\"\"></section><section class=\"my-8\"><input type=\"submit\" name=\"submit\" value=\"Submit\" class=\"bg-gray-400 px-4 py-2 rounded-full mx-auto block hover:cursor-pointer\"></section></form>")
			if err != nil {
				return err
			}
			err = UnsafeRawHtml(f.Errors["overall"]).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</section>")
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
