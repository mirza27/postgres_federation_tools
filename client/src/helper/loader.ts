export interface DefaultLoaderResponse<T> {
  data: T | null;
  ok: boolean;
  message: string;
}
