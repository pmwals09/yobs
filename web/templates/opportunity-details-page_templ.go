// Code generated by templ@v0.2.364 DO NOT EDIT.

package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"strconv"

	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/user"
)

func OpportunityDetailsPage(user *user.User, od helpers.OpptyDetails, userDocuments []document.Document, fd helpers.FormData) templ.Component {
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
			var var_3 string = od.Oppty.CompanyName
			_, err = templBuffer.WriteString(templ.EscapeString(var_3))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" ")
			if err != nil {
				return err
			}
			var_4 := `- `
			_, err = templBuffer.WriteString(var_4)
			if err != nil {
				return err
			}
			var var_5 string = od.Oppty.Role
			_, err = templBuffer.WriteString(templ.EscapeString(var_5))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h1> <section class=\"mb-4\"><dl class=\"grid grid-cols-[min-content_auto] gap-x-4\"><dt class=\"font-bold whitespace-nowrap\">")
			if err != nil {
				return err
			}
			var_6 := `Company Name:`
			_, err = templBuffer.WriteString(var_6)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</dt><dd>")
			if err != nil {
				return err
			}
			var var_7 string = od.Oppty.CompanyName
			_, err = templBuffer.WriteString(templ.EscapeString(var_7))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</dd><dt class=\"font-bold whitespace-nowrap\">")
			if err != nil {
				return err
			}
			var_8 := `Role:`
			_, err = templBuffer.WriteString(var_8)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</dt><dd>")
			if err != nil {
				return err
			}
			var var_9 string = od.Oppty.Role
			_, err = templBuffer.WriteString(templ.EscapeString(var_9))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</dd><dt class=\"font-bold whitespace-nowrap\">")
			if err != nil {
				return err
			}
			var_10 := `URL:`
			_, err = templBuffer.WriteString(var_10)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</dt><dd>")
			if err != nil {
				return err
			}
			var var_11 string = od.Oppty.URL
			_, err = templBuffer.WriteString(templ.EscapeString(var_11))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</dd><dt class=\"font-bold whitespace-nowrap\">")
			if err != nil {
				return err
			}
			var_12 := `Status:`
			_, err = templBuffer.WriteString(var_12)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</dt><dd>")
			if err != nil {
				return err
			}
			var var_13 string = string(od.Oppty.Status)
			_, err = templBuffer.WriteString(templ.EscapeString(var_13))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</dd><dt class=\"font-bold whitespace-nowrap\">")
			if err != nil {
				return err
			}
			var_14 := `Application Date:`
			_, err = templBuffer.WriteString(var_14)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</dt><dd>")
			if err != nil {
				return err
			}
			var var_15 string = formatApplicationDate(od.Oppty.ApplicationDate)
			_, err = templBuffer.WriteString(templ.EscapeString(var_15))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</dd></dl></section> <section class=\"mb-4\"><h2 class=\"text-[1.5em]\">")
			if err != nil {
				return err
			}
			var_16 := `Contacts`
			_, err = templBuffer.WriteString(var_16)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h2><table class=\"table-auto border-collapse w-full mb-4\"><thead></thead><tbody>")
			if err != nil {
				return err
			}
			for _, contact := range od.Contacts {
				_, err = templBuffer.WriteString("<tr><td class=\"border-b\">")
				if err != nil {
					return err
				}
				var var_17 string = contact.Name
				_, err = templBuffer.WriteString(templ.EscapeString(var_17))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</td><td class=\"border-b\">")
				if err != nil {
					return err
				}
				var var_18 string = contact.Company
				_, err = templBuffer.WriteString(templ.EscapeString(var_18))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</td><td class=\"border-b\">")
				if err != nil {
					return err
				}
				var var_19 string = contact.Title
				_, err = templBuffer.WriteString(templ.EscapeString(var_19))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</td><td class=\"border-b\">")
				if err != nil {
					return err
				}
				var var_20 string = contact.Phone
				_, err = templBuffer.WriteString(templ.EscapeString(var_20))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</td><td class=\"border-b\">")
				if err != nil {
					return err
				}
				var var_21 string = contact.Email
				_, err = templBuffer.WriteString(templ.EscapeString(var_21))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</td></tr>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</tbody></table></section> <section class=\"mb-4\"><h2 class=\"text-[1.5em]\">")
			if err != nil {
				return err
			}
			var_22 := `Job Description`
			_, err = templBuffer.WriteString(var_22)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h2> ")
			if err != nil {
				return err
			}
			var var_23 string = od.Oppty.Description
			_, err = templBuffer.WriteString(templ.EscapeString(var_23))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</section> <section class=\"mb-4\"><h2 class=\"text-[1.5em]\">")
			if err != nil {
				return err
			}
			var_24 := `Tasks`
			_, err = templBuffer.WriteString(var_24)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h2></section> <section id=\"attachments-section\" class=\"mb-4\"><h2 class=\"text-[1.5em]\">")
			if err != nil {
				return err
			}
			var_25 := `Attachments`
			_, err = templBuffer.WriteString(var_25)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h2>")
			if err != nil {
				return err
			}
			if od.Documents != nil && len(od.Documents) > 0 {
				_, err = templBuffer.WriteString("<table class=\"table-auto border-collapse w-full mb-4\"><thead><tr><th class=\"border-b text-left\">")
				if err != nil {
					return err
				}
				var_26 := `Title`
				_, err = templBuffer.WriteString(var_26)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</th><th class=\"border-b text-left\">")
				if err != nil {
					return err
				}
				var_27 := `Type`
				_, err = templBuffer.WriteString(var_27)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</th><th class=\"border-b text-left\">")
				if err != nil {
					return err
				}
				var_28 := `File Name`
				_, err = templBuffer.WriteString(var_28)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</th></tr></thead><tbody>")
				if err != nil {
					return err
				}
				for _, doc := range od.Documents {
					_, err = templBuffer.WriteString("<tr class=\"hover:bg-sky-100\"><td class=\"border-b\">")
					if err != nil {
						return err
					}
					var var_29 string = doc.Title
					_, err = templBuffer.WriteString(templ.EscapeString(var_29))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</td><td class=\"border-b\">")
					if err != nil {
						return err
					}
					var var_30 string = string(doc.Type)
					_, err = templBuffer.WriteString(templ.EscapeString(var_30))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</td><td class=\"border-b\"><a href=\"")
					if err != nil {
						return err
					}
					var var_31 templ.SafeURL = templ.URL(doc.URL)
					_, err = templBuffer.WriteString(templ.EscapeString(string(var_31)))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\">")
					if err != nil {
						return err
					}
					var var_32 string = doc.FileName
					_, err = templBuffer.WriteString(templ.EscapeString(var_32))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</a></td></tr>")
					if err != nil {
						return err
					}
				}
				_, err = templBuffer.WriteString("</tbody></table>")
				if err != nil {
					return err
				}
			}
			err = RenderStandardError(fd.Errors["document-table"]).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			if userDocuments != nil && len(userDocuments) > 0 {
				_, err = templBuffer.WriteString("<form method=\"POST\" action=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(insertIDIntoString("/opportunities/{}/attach-existing", od.Oppty.ID)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\" hx-post=\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(insertIDIntoString("/opportunities/{}/attach-existing", od.Oppty.ID)))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\" hx-target=\"#attachments-section\" class=\"w-1/2\"><section class=\"grid grid-cols-2 gap-y-2\"><label for=\"existing-attachment\">")
				if err != nil {
					return err
				}
				var_33 := `Existing Attachments`
				_, err = templBuffer.WriteString(var_33)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</label><div><select name=\"existing-attachment\"><option value=\"\">")
				if err != nil {
					return err
				}
				var_34 := `Select a document...`
				_, err = templBuffer.WriteString(var_34)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option>")
				if err != nil {
					return err
				}
				for _, doc := range userDocuments {
					_, err = templBuffer.WriteString("<option value=\"")
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(templ.EscapeString(strconv.Itoa(int(doc.ID))))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("\">")
					if err != nil {
						return err
					}
					var var_35 string = doc.Title
					_, err = templBuffer.WriteString(templ.EscapeString(var_35))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString(" ")
					if err != nil {
						return err
					}
					var_36 := `- `
					_, err = templBuffer.WriteString(var_36)
					if err != nil {
						return err
					}
					var var_37 string = doc.FileName
					_, err = templBuffer.WriteString(templ.EscapeString(var_37))
					if err != nil {
						return err
					}
					_, err = templBuffer.WriteString("</option>")
					if err != nil {
						return err
					}
				}
				_, err = templBuffer.WriteString("</select>")
				if err != nil {
					return err
				}
				err = RenderStandardError(fd.Errors["existing-attachment"]).Render(ctx, templBuffer)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div></section><input type=\"submit\" name=\"attachment-submit\" value=\"Add attached\" class=\"bg-gray-400 px-4 py-2 rounded-full mx-auto block hover:cursor-pointer my-3\"></form>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("<form method=\"POST\" action=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(insertIDIntoString("/opportunities/{}/upload", od.Oppty.ID)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-post=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(insertIDIntoString("/opportunities/{}/upload", od.Oppty.ID)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\" hx-target=\"#attachments-section\" hx-select=\"#attachments-section\" enctype=\"multipart/form-data\" class=\"w-1/2\"><section class=\"grid grid-cols-2 gap-y-2\"><label for=\"attachment-name\">")
			if err != nil {
				return err
			}
			var_38 := `Attachment Name`
			_, err = templBuffer.WriteString(var_38)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><div><input type=\"text\" name=\"attachment-name\" value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(fd.Values["attachment-name"]))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			err = RenderStandardError(fd.Errors["attachment-name"]).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div><label for=\"attachment-type\">")
			if err != nil {
				return err
			}
			var_39 := `Attachment Type`
			_, err = templBuffer.WriteString(var_39)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><div><select name=\"attachment-type\" value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(fd.Values["attachment-type"]))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\"><option value=\"Resume\">")
			if err != nil {
				return err
			}
			var_40 := `Resume`
			_, err = templBuffer.WriteString(var_40)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</option><option value=\"CoverLetter\">")
			if err != nil {
				return err
			}
			var_41 := `Cover Letter`
			_, err = templBuffer.WriteString(var_41)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</option><option value=\"Communication\">")
			if err != nil {
				return err
			}
			var_42 := `Communication`
			_, err = templBuffer.WriteString(var_42)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</option></select>")
			if err != nil {
				return err
			}
			err = RenderStandardError(fd.Errors["attachment-type"]).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div><label for=\"attachment-file\">")
			if err != nil {
				return err
			}
			var_43 := `PDF Attachment`
			_, err = templBuffer.WriteString(var_43)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label><div><input type=\"file\" name=\"attachment-file\" accept=\".pdf\" value=\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(fd.Values["attachment-file"]))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\">")
			if err != nil {
				return err
			}
			err = RenderStandardError(fd.Errors["attachment-file"]).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div></section><input type=\"submit\" name=\"attachment-submit\" value=\"Submit new\" class=\"bg-gray-400 px-4 py-2 rounded-full mx-auto block hover:cursor-pointer mt-3\"></form></section>")
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
