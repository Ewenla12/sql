package routes

import (
	"sql/handlers"
	"sql/middleware"

	"github.com/gofiber/fiber/v3"
)

// adminRoles lists every role in your admin_role enum — used wherever an
// endpoint should be reachable by any kind of admin, not students.
var adminRoles = []string{"hod", "faculty_officer", "dean", "registrar", "super_admin"}

// RegisterRoutes attaches all API routes to the given Fiber app.
func RegisterRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/student/login", handlers.StudentLoginHandler)
	auth.Post("/admin/login", handlers.AdminLoginHandler)

	students := app.Group("/students", middleware.RequireAuth)

	students.Get("/", middleware.RequireRole(adminRoles...), handlers.GetStudentsHandler)
	students.Get("/:student_id", middleware.RequireRole(adminRoles...), handlers.GetStudentByIDHandler)
	students.Post("/", middleware.RequireRole(adminRoles...), handlers.CreateStudentHandler)
	students.Put("/:student_id", middleware.RequireRole(adminRoles...), handlers.UpdateStudentHandler)
	students.Delete("/:student_id", middleware.RequireRole(adminRoles...), handlers.DeleteStudentHandler)

	students.Get("/me/courses", handlers.GetMyCoursesHandler)
	students.Post("/me/courses/register", handlers.ApplyForCoursesHandler)

	admin := app.Group("/admin", middleware.RequireAuth, middleware.RequireRole(adminRoles...))
	admin.Get("/course-registrations", handlers.GetCourseRegistrationsHandler)
	admin.Post("/course-registrations/:registration_id/approve", handlers.ApproveCourseRegistrationHandler)
	admin.Post("/course-registrations/:registration_id/reject", handlers.RejectCourseRegistrationHandler)

	app.Get("/faculties", handlers.GetFacultiesHandler)
	app.Get("/departments", handlers.GetDepartmentsHandler)
	app.Get("/courses", handlers.GetCoursesHandler)
}
