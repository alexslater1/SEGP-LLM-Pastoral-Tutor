export enum UserRoleEnum {
  ADMIN = "admin",
  STUDENT = "student",
  TUTOR = "tutor",
}

export const ADMIN_PANEL_ROLES: UserRoleEnum[] = [UserRoleEnum.ADMIN, UserRoleEnum.TUTOR];

export const PUBLIC_PAGES = [/\/sign-in.*/, /\/sign-up.*/, /\/fonts.*/, /\/auth.*/, /\/$/];

export const STUDENT_WHITELISTED_PAGES = [/\/chat.*/].concat(PUBLIC_PAGES);
export const TUTOR_WHITELISTED_PAGES = [/\/admin$/, /\/admin\/tutees.*/, /\/admin\/tutee-chat.*/].concat(STUDENT_WHITELISTED_PAGES);
export const ADMIN_WHITELISTED_PAGES = [/.*/];
