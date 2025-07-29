export interface Job {
  filename: string;
  status: string;
  type: string;
  startedAt: string;
  completedAt: string;
  output: string;
  error: string;
  createdAt: string;
  updatedAt: string;
}

export class Job implements Job {
  id: string | null = null;
  filename = "";
  status = "";
  type = "";
  startedAt = "";
  completedAt = "";
  output = "";
  error = "";
  createdAt = "";
  updatedAt = "";
}
