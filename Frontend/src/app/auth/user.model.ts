import { Profile } from "./profile.model";

export interface User {
  username: string;
  password: string;
  mail: string;
  role: string;
  blocked: boolean;
  profile: Profile;
}
