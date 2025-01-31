"use client";
import useSWR from "swr";
import { Post } from "../models/post";

type ApiResponse = Post[];

type ApiError = {
  message: string;
  status: number;
};

export const pollFollowingPosts = () => {
  const POLLING_TIMEOUT = 30000;

  const fetcher = async (url: string) => {
    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), POLLING_TIMEOUT);

      const response = await fetch(url, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
        signal: controller.signal,
      });

      clearTimeout(timeoutId);

      if (
        response.status === 204 ||
        (error instanceof Error && error.name === "AbortError")
      ) {
        return null;
      }

      if (!response.ok) {
        const error: ApiError = {
          message: `Server error: ${response.statusText}`,
          status: response.status,
        };
        throw error;
      }

      return await response.json();
    } catch (error) {
      if (error instanceof Error && error.name === "AbortError") {
        return null;
      }
      throw error;
    }
  };

  const { data, error, mutate } = useSWR<ApiResponse | null, ApiError>(
    `${process.env.NEXT_PUBLIC_LOCAL_API_BASE_URL}/api/${process.env.NEXT_PUBLIC_USER_ID}/polling?event_type=PollingRequest`,
    fetcher,
    {
      refreshWhenHidden: true,
      revalidateOnFocus: false,
      revalidateOnReconnect: true,
      onSuccess: () => {
        mutate();
      },
    }
  );

  return {
    data: data || undefined,
    error,
    refresh: () => mutate(),
  };
};
