import { create } from 'zustand';

export type MessageRole = 'user' | 'agent' | 'system';

export interface ChatMessage {
  id: string;
  role: MessageRole;
  content: string;
  traceId?: string;
  timestamp: number;
}

interface ChatState {
  messages: ChatMessage[];
  isTyping: boolean;
  addMessage: (message: ChatMessage) => void;
  setTyping: (typing: boolean) => void;
  appendContent: (id: string, delta: string) => void;
}

export const useChat = create<ChatState>((set) => ({
  messages: [
    {
      id: 'welcome',
      role: 'agent',
      content: 'Ask me to find today\'s highest-probability swing trades. I\'ll coordinate the agents.',
      timestamp: Date.now(),
    },
  ],
  isTyping: false,
  addMessage: (message) => set((s) => ({ messages: [...s.messages, message] })),
  setTyping: (typing) => set({ isTyping: typing }),
  appendContent: (id, delta) => set((s) => ({
    messages: s.messages.map((m) =>
      m.id === id ? { ...m, content: m.content + delta } : m
    ),
  })),
}));
