package handlers

import (
	"strconv"

	"sql/repository"
	"sql/utils"

	"github.com/gofiber/fiber/v3"
)

// getClaims is a small helper to pull the authenticated user's claims
// out of the request context (set earlier by middleware.RequireAuth).
func getClaims(c fiber.Ctx) (*utils.Claims, bool) {
	claims, ok := c.Locals("claims").(*utils.Claims)
	return claims, ok
}

// ApplyForCoursesHandler handles POST /students/me/courses/register
// The student applying is taken from the JWT, not from the request body —
// a student can only register courses for themselves.
func ApplyForCoursesHandler(c fiber.Ctx) error {
	claims, ok := getClaims(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "not authenticated"})
	}

	type request struct {
		AcademicSession string `json:"academic_session"`
		CourseIDs       []int  `json:"course_ids"`
	}

	var req request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	ids, err := repository.ApplyForCourses(c.Context(), claims.Subject, req.AcademicSession, req.CourseIDs)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"registration_ids": ids,
		"status":           "pending",
		"message":          "Course registration submitted for approval",
	})
}

// GetMyCoursesHandler handles GET /students/me/courses
func GetMyCoursesHandler(c fiber.Ctx) error {
	claims, ok := getClaims(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "not authenticated"})
	}

	session := c.Query("academic_session")

	courses, err := repository.GetStudentCourses(c.Context(), claims.Subject, session)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"student_id":       claims.Subject,
		"academic_session": session,
		"courses":          courses,
	})
}

// GetCourseRegistrationsHandler handles GET /admin/course-registrations
func GetCourseRegistrationsHandler(c fiber.Ctx) error {
	filter := repository.CourseRegistrationFilter{
		Status:          c.Query("status", "pending"),
		AcademicSession: c.Query("academic_session"),
	}

	if deptID, err := strconv.Atoi(c.Query("department_id")); err == nil {
		filter.DepartmentID = deptID
	}

	regs, err := repository.GetCourseRegistrations(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": regs})
}

// ApproveCourseRegistrationHandler handles POST /admin/course-registrations/:registration_id/approve
func ApproveCourseRegistrationHandler(c fiber.Ctx) error {
	claims, ok := getClaims(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "not authenticated"})
	}

	regID, err := strconv.Atoi(c.Params("registration_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid registration_id"})
	}

	adminID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid admin token"})
	}

	if err := repository.ApproveRegistration(c.Context(), regID, adminID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"registration_id": regID,
		"status":          "approved",
		"message":         "Course registration approved",
	})
}

// RejectCourseRegistrationHandler handles POST /admin/course-registrations/:registration_id/reject
func RejectCourseRegistrationHandler(c fiber.Ctx) error {
	claims, ok := getClaims(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "not authenticated"})
	}

	regID, err := strconv.Atoi(c.Params("registration_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid registration_id"})
	}

	adminID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid admin token"})
	}

	if err := repository.RejectRegistration(c.Context(), regID, adminID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"registration_id": regID,
		"status":          "rejected",
		"message":         "Course registration rejected",
	})
}
