export interface User {
  email: string;
  password: string;
}

export class User implements User {
  id: string | null = null;
  email = "";
  password = "";
  token: string | null = null;
  role: string | null = null;
  created_at: string | null = null;
  updated_at: string | null = null;
}
